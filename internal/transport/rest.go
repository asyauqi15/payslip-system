package transport

import (
	"context"
	"fmt"
	"github.com/asyauqi15/payslip-system/internal"
	"github.com/asyauqi15/payslip-system/internal/transport/middleware"
	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
	"net/http"
)

type RESTServer struct {
	srv *http.Server
}

func NewRESTServer(
	config internal.Config,
	db *gorm.DB,
) (*RESTServer, error) {
	//jwtAuth, err := jwtauth.NewJWTAuthentication(config.HTTPServer)
	//if err != nil {
	//	return nil, err
	//}

	routes := chi.NewRouter()
	routes.Use(middleware.Recoverer)
	routes.Use(middleware.CheckIPAddress)

	//routes.Route("/v1", func(v1 chi.Router) {
	//	v1.Use(chiMiddleware.StripSlashes)
	//	v1.Mount("/auth", endpoint.NewAuthEndpoint(service.NewAuthService(jwtAuth, s.UserSvc, r.RefreshTokenRepository)))
	//	v1.Mount("/activation", endpoint.NewActivationEndpoint(r.ApplicationConfigurationRepo, s.ActivationSvc))
	//	v1.Get("/bulk-invoices/{bulkInvoiceID}/pay", bulkInvoiceEndpoint.Pay)
	//	v1.Group(func(transactions chi.Router) {
	//		transactions.Mount("/transactions", endpoint.NewTransactionEndpoint(s.TransactionSvc, s.UserSvc))
	//	})
	//	v1.Group(func(mekariCapital chi.Router) {
	//		mekariCapital.Use(authz.NewMekariCapitalAuthz(config.MekariCapital.APIKey).Authorizer)
	//		mekariCapital.Mount("/mekari-capitals", endpoint.NewMekariCapitalEndpoint(s.MekariCapitalSvc))
	//	})
	//
	//	v1.Group(func(webhook chi.Router) {
	//		webhook.Use(authz.NewWebhookAuthz(config.Webhook.APIKey).Authorizer)
	//		webhook.Mount("/webhooks", endpoint.NewWebhookEndpoint(s.WebhookSvc))
	//	})
	//
	//	v1.Mount("/public", endpoint.NewPublicEndpoint(s.JurnalSvc, s.PreactivationEmailSvc, s.InvoicePageSvc, config))
	//
	//	v1.Group(func(protected chi.Router) {
	//		protected.Use(jwtAuth.Authenticator, authz.CasbinAuthorizer, jurnalAccessTokenResolver.Resolver, middleware.SentryIdentifyUserMiddleware)
	//		protected.Mount("/users", endpoint.NewUserEndpoint(s.UserSvc))
	//		protected.Mount("/settings", endpoint.NewSettingEndpoint(s.SettingSvc, r.BankDetailRepo, r.UserAuditRepo, r.PaymentMethodConfigurationRepo, config.Application))
	//		protected.Mount("/companies", endpoint.NewCompanyAccountEndpoint(s.CompanyAccountSvc))
	//		protected.Mount("/customers", endpoint.NewCustomerEndpoint(s.CustomerSvc))
	//		protected.Mount("/bulk-invoices", bulkInvoiceEndpoint.ProtectedEndpoints())
	//		protected.Mount("/fwd-kyc", endpoint.NewKYCEndpoint(s.KYCSvc))
	//		protected.Mount("/kyc-generals", endpoint.NewKYCGeneralEndpoint(s.KYCGeneralSvc))
	//		protected.Mount("/dashboard", endpoint.NewDashboardEndpoint(s.DashboardSvc))
	//		protected.Mount("/onboardings", endpoint.NewOnboardingFlagEndpoint(s.OnboardingFlagSvc))
	//		protected.Mount("/sales-invoices", salesInvoiceEndpoint.Endpoints())
	//		protected.Mount("/fwd-quota", endpoint.NewQuotaEndpoint(s.QuotaSvc))
	//		protected.Mount("/sales-orders", endpoint.NewSalesOrderEndpoint(s.SalesOrderSvc, config))
	//		protected.Mount("/proforma-sales-invoices", proformaSalesInvoiceEndpoint.Endpoints())
	//		protected.Mount("/proforma-orders", endpoint.NewProformaOrderEndpoint(s.ProformaOrderSvc))
	//		protected.Mount("/jurnal", endpoint.NewJurnalEndpoint(r.JurnalSdkRepo, s.JurnalSvc, config.Feature))
	//		protected.Mount("/mekari-payment", endpoint.NewMekariPaymentEndpoint(r.BankListRepo, s.BankDetailSvc))
	//		protected.Mount("/withdrawal", endpoint.NewWithdrawalEndpoint(r.WithdrawalRepo, s.WithdrawalSvc))
	//		protected.Route("/payouts", func(payouts chi.Router) {
	//			payouts.Mount("/", payoutEndpoint.NewWebEndpointHandler(db, r.JurnalSdkRepo, config))
	//			payouts.Mount("/pending-approval", endpoint.NewPayoutEndpoint(s.PayoutSvc, config))
	//			payouts.Mount("/payment-histories", payoutEndpoint.NewWebPaymentHistoryEndpoint(db, r.JurnalSdkRepo, config))
	//			payouts.Mount("/payment-requests", payoutRequestEndpoint.NewWebEndpoint(db, r.JurnalSdkRepo, r.MekariPaymentRepo, worker, locker, config, s.JurnalSvc, s.PayoutValidationSvc))
	//		})
	//		protected.Route("/capital", func(capi chi.Router) {
	//			capi.Mount("/banner", endpoint.NewCapitalBannerEndpoint(r.CapitalRepo, config.MekariCapital))
	//		})
	//		protected.Mount("/balances/transaction-histories", transactionHistoryEndoint.NewWebEndpoint(worker, db, r.JurnalSdkRepo))
	//		protected.Mount("/balances/withdrawals", withdrawalEndpoint.NewWebBalanceEndpoint(db))
	//		protected.Mount("/topup", topupEndpoint.NewWebEndpoint(db, r.MekariPaymentRepo, r.JurnalSdkRepo, s.KYCSvc, config))
	//	})
	//	v1.Group(func(internalGroup chi.Router) {
	//		internalGroup.Use(authz.NewInternalAuthz(config.Internal.APIKey).Authorizer)
	//		internalGroup.Route("/internals", func(internalRoute chi.Router) {
	//			internalRoute.Mount("/", endpoint.NewInternalEndpoint(
	//				s.BankDetailSvc, s.InternalSvc, s.PreactivationEmailSvc,
	//				r.SalesInvoicePaymentRepo, s.ActivationSvc, r.PurchasePaymentRepo, r.LedgerRepo, r.ExpensePaymentRepo))
	//
	//			internalRoute.Mount("/payouts", payoutEndpoint.NewInternalEndpoint(db, worker, worker, config, r.OnesignalRepo, r.JurnalWebRepo))
	//			internalRoute.Mount("/payments", paymentEndpoint.NewInternalEndpoint(db, worker, r.JurnalSdkRepo, config, r.OnesignalRepo, r.JurnalWebRepo, r.BackyardRepo))
	//			internalRoute.Mount("/sso", mekariSsoEndpoint.NewInternalEndpoint(worker))
	//		})
	//		internalGroup.Mount("/balances", endpoint.NewBalanceEndpoint(s.UserBalanceSvc))
	//		internalGroup.Mount("/backyard", endpoint.NewBackyardEndpoint(r.CacheRepo))
	//	})
	//
	//	v1.Group(func(internal chi.Router) {
	//		internal.Use(authz.NewMekariPaymentCallbackAuthz(config.MekariPayment).Authorizer)
	//		internal.Mount("/mekari-payment/callback", endpoint.NewCallbackEndpoint(s.CallbackSvc, r.PayoutRepo, worker))
	//	})
	//
	//	v1.Route("/mobile", func(mobile chi.Router) {
	//		mobile.Use(authz.NewInternalAuthz(config.Mobile.APIKey).Authorizer, middleware.IdentifySource)
	//
	//		mobile.Mount("/activation", activation.NewMobileEndpoint(db, r.JurnalSdkRepo))
	//		mobile.Mount("/activation-config", contactus.NewMobileEndpoint(db))
	//
	//		mobile.Group(func(jurnalAuth chi.Router) {
	//			jurnalAuth.Use(auth.NewJurnalTokenAuthentication(s.JurnalAuthSvc).Authenticator)
	//
	//			jurnalAuth.Mount("/role-info", role.NewMobileEndpoint(db))
	//			jurnalAuth.Mount("/balance", balance.NewMobileEndpoint(db, r.MekariPaymentRepo))
	//			jurnalAuth.Mount("/payouts", payoutEndpoint.NewMobileEndpoint(db, worker, r.MekariPaymentRepo, r.JurnalSdkRepo, config, locker, s.PayoutValidationSvc))
	//			jurnalAuth.Mount("/fee-config", feeconfig.NewMobileEndpoint(db))
	//			jurnalAuth.Mount("/jurnal", jurnalEndpoint.NewMobileEndpoint(r.JurnalSdkRepo))
	//		})
	//	})
	//})

	routes.Get("/ping", pingHandler)
	swaggerRoutes(routes)

	return &RESTServer{
		srv: &http.Server{
			Addr:              fmt.Sprintf("0.0.0.0:%d", config.HTTPServer.Port),
			Handler:           routes,
			ReadTimeout:       config.HTTPServer.ReadTimeout,
			ReadHeaderTimeout: config.HTTPServer.ReadHeaderTimeout,
			IdleTimeout:       config.HTTPServer.IdleTimeout,
			WriteTimeout:      config.HTTPServer.WriteTimeout,
		},
	}, nil
}

func (r *RESTServer) Start() error {
	return r.srv.ListenAndServe()
}

func (r *RESTServer) Stop(ctx context.Context) error {
	return r.srv.Shutdown(ctx)
}
