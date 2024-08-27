package controller

import (
	"net/http"
	"time"

	"github.com/dagota12/Loan-Tracker/bootstrap"
	"github.com/dagota12/Loan-Tracker/domain"
	"github.com/dagota12/Loan-Tracker/internal/emailutil"
	"github.com/dagota12/Loan-Tracker/internal/otputil"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type ResetPasswordController struct {
	ResetPasswordUsecase domain.ResetPasswordUsecase
	Env                  *bootstrap.Env
}

func NewResetPasswordController(env *bootstrap.Env, resetPasswordUsecase domain.ResetPasswordUsecase) *ResetPasswordController {
	return &ResetPasswordController{
		ResetPasswordUsecase: resetPasswordUsecase,
		Env:                  env,
	}
}

func (rc *ResetPasswordController) ForgotPassword(ctx *gin.Context) {
	var req domain.ForgotPasswordRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	_, err := rc.ResetPasswordUsecase.GetUserByEmail(ctx, req.Email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if _, err := rc.ResetPasswordUsecase.GetOTPByEmail(ctx, req.Email); err == nil {
		err = rc.ResetPasswordUsecase.DeleteOtp(ctx, req.Email)
		if err != nil {
			ctx.Error(err)
			return
		}
	}
	otp, err := otputil.GenerateOTP()
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	err = emailutil.SendOtpVerificationEmail(req.Email, otp, rc.Env)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	hashedcode, err := bcrypt.GenerateFromPassword([]byte(otp), bcrypt.DefaultCost)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	newOtp := domain.OtpSave{
		Email:     req.Email,
		Code:      string(hashedcode),
		ExpiresAt: time.Now().Add(time.Minute * time.Duration(rc.Env.PassResetCodeExpirationMin)),
	}
	err = rc.ResetPasswordUsecase.SaveOtp(ctx, &newOtp)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "OTP sent successfully"})
}

func (rc *ResetPasswordController) ResetPassword(ctx *gin.Context) {
	var req domain.ResetPasswordRequest

	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	user, err := rc.ResetPasswordUsecase.GetUserByEmail(ctx, req.Email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	originalOtp, err := rc.ResetPasswordUsecase.GetOTPByEmail(ctx, req.Email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if time.Now().After(originalOtp.ExpiresAt) {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "OTP expired"})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(originalOtp.Code), []byte(req.Code))
	if err != nil {
		ctx.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid OTP"})
		return
	}

	err = rc.ResetPasswordUsecase.ResetPassword(ctx, user.ID.Hex(), &req)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	err = rc.ResetPasswordUsecase.DeleteOtp(ctx, req.Email)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{"message": "Password reset successfully"})
}
