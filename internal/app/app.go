package app

import (
	"context"
	"log"
	"net/http"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	httpSwagger "github.com/swaggo/http-swagger"
	"go.uber.org/fx"

	"github.com/viniciuscluna/tc-fiap-50/internal/infrastructure/clients"
	"github.com/viniciuscluna/tc-fiap-50/internal/shared/config"
	"github.com/viniciuscluna/tc-fiap-50/internal/shared/httpclient"

	orderController "github.com/viniciuscluna/tc-fiap-50/internal/order/controller"
	orderRepositories "github.com/viniciuscluna/tc-fiap-50/internal/order/domain/repositories"
	orderApiController "github.com/viniciuscluna/tc-fiap-50/internal/order/infrastructure/api/controller"
	orderPersistence "github.com/viniciuscluna/tc-fiap-50/internal/order/infrastructure/persistence"
	orderPresenter "github.com/viniciuscluna/tc-fiap-50/internal/order/presenter"
	orderUseCasesAdd "github.com/viniciuscluna/tc-fiap-50/internal/order/usecase/addOrder"
	orderUseCasesGet "github.com/viniciuscluna/tc-fiap-50/internal/order/usecase/getOrder"
	orderUseCasesGetOrderStatus "github.com/viniciuscluna/tc-fiap-50/internal/order/usecase/getOrderStatus"
	orderUseCasesGetOrders "github.com/viniciuscluna/tc-fiap-50/internal/order/usecase/getOrders"
	orderUseCasesUpdateOrderStatus "github.com/viniciuscluna/tc-fiap-50/internal/order/usecase/updateOrderStatus"

	"github.com/viniciuscluna/tc-fiap-50/pkg/rest"
	"github.com/viniciuscluna/tc-fiap-50/pkg/storage/postgres"
)

func InitializeApp() *fx.App {
	return fx.New(
		fx.Provide(
			// Configuration
			config.Load,

			// Database
			postgres.NewPostgresDB,

			// HTTP Client
			func(cfg *config.Config) httpclient.HTTPClient {
				return httpclient.NewHTTPClient(cfg.HTTPClientTimeout, cfg.HTTPClientRetryCount, cfg.HTTPClientRetryBackoff)
			},

			// External Service Clients
			fx.Annotate(
				func(httpClient httpclient.HTTPClient, cfg *config.Config) clients.CustomerClient {
					return clients.NewCustomerClientImpl(httpClient, cfg.CustomerServiceURL)
				},
				fx.As(new(clients.CustomerClient)),
			),
			fx.Annotate(
				func(httpClient httpclient.HTTPClient, cfg *config.Config) clients.ProductClient {
					return clients.NewProductClientImpl(httpClient, cfg.ProductServiceURL)
				},
				fx.As(new(clients.ProductClient)),
			),

			// Order Repositories
			fx.Annotate(orderPersistence.NewOrderRepositoryImpl, fx.As(new(orderRepositories.OrderRepository))),
			fx.Annotate(orderPersistence.NewOrderProductRepositoryImpl, fx.As(new(orderRepositories.OrderProductRepository))),
			fx.Annotate(orderPersistence.NewOrderStatusRepositoryImpl, fx.As(new(orderRepositories.OrderStatusRepository))),

			// Order Use Cases (now with client dependencies)
			fx.Annotate(orderUseCasesAdd.NewAddOrderUseCaseImpl, fx.As(new(orderUseCasesAdd.AddOrderUseCase))),
			fx.Annotate(orderUseCasesGet.NewGetOrderUseCaseImpl, fx.As(new(orderUseCasesGet.GetOrderUseCase))),
			fx.Annotate(orderUseCasesGetOrders.NewGetOrdersUseCaseImpl, fx.As(new(orderUseCasesGetOrders.GetOrdersUseCase))),
			fx.Annotate(orderUseCasesGetOrderStatus.NewGetOrderStatusUseCaseImpl, fx.As(new(orderUseCasesGetOrderStatus.GetOrderStatusUseCase))),
			fx.Annotate(orderUseCasesUpdateOrderStatus.NewUpdateOrderStatusUseCaseImpl, fx.As(new(orderUseCasesUpdateOrderStatus.UpdateOrderStatusUseCase))),

			// Order Controller and Presenter (with client dependencies)
			fx.Annotate(orderController.NewOrderControllerImpl, fx.As(new(orderController.OrderController))),
			fx.Annotate(orderPresenter.NewOrderPresenterImpl, fx.As(new(orderPresenter.OrderPresenter))),
			chi.NewRouter,
			func(orderController orderController.OrderController) []rest.Controller {
				return []rest.Controller{
					orderApiController.NewOrderController(orderController),
				}
			},
		),
		fx.Invoke(registerRoutes),
		fx.Invoke(startHTTPServer),
	)
}

func registerRoutes(r *chi.Mux, controllers []rest.Controller) {
	r.Use(middleware.Logger)

	// Swagger UI
	r.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"), // The URL pointing to API definition
	))

	for _, controller := range controllers {
		controller.RegisterRoutes(r)
	}
}

func startHTTPServer(lc fx.Lifecycle, r *chi.Mux) {
	lc.Append(fx.Hook{
		OnStart: func(ctx context.Context) error {
			go func() {
				log.Println("Starting HTTP server on :8080")
				if err := http.ListenAndServe(":8080", r); err != nil {
					log.Fatalf("Failed to start HTTP server: %v", err)
				}
			}()
			return nil
		},
		OnStop: func(ctx context.Context) error {
			log.Println("Shutting down HTTP server gracefully")
			return nil
		},
	})
}
