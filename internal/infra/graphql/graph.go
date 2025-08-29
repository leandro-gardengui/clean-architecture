package graphql

import (
	"clean-architecture/internal/usecase"

	"github.com/graphql-go/graphql"
)

type Graph struct {
	CreateOrderUseCase *usecase.CreateOrderUseCase
	ListOrdersUseCase  *usecase.ListOrdersUseCase
}

func NewGraph(createOrderUseCase *usecase.CreateOrderUseCase, listOrdersUseCase *usecase.ListOrdersUseCase) *Graph {
	return &Graph{
		CreateOrderUseCase: createOrderUseCase,
		ListOrdersUseCase:  listOrdersUseCase,
	}
}

var orderType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Order",
	Fields: graphql.Fields{
		"id": &graphql.Field{
			Type: graphql.String,
		},
		"price": &graphql.Field{
			Type: graphql.Float,
		},
		"tax": &graphql.Field{
			Type: graphql.Float,
		},
		"final_price": &graphql.Field{
			Type: graphql.Float,
		},
	},
})

func (g *Graph) QueryType() *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"listOrders": &graphql.Field{
				Type: graphql.NewList(orderType),
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					orders, err := g.ListOrdersUseCase.Execute()
					if err != nil {
						return nil, err
					}

					var result []map[string]interface{}
					for _, order := range orders {
						result = append(result, map[string]interface{}{
							"id":          order.ID,
							"price":       order.Price,
							"tax":         order.Tax,
							"final_price": order.FinalPrice,
						})
					}
					return result, nil
				},
			},
		},
	})
}

func (g *Graph) MutationType() *graphql.Object {
	return graphql.NewObject(graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
			"createOrder": &graphql.Field{
				Type: orderType,
				Args: graphql.FieldConfigArgument{
					"price": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Float),
					},
					"tax": &graphql.ArgumentConfig{
						Type: graphql.NewNonNull(graphql.Float),
					},
				},
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					price := p.Args["price"].(float64)
					tax := p.Args["tax"].(float64)

					input := usecase.CreateOrderInputDTO{
						Price: price,
						Tax:   tax,
					}

					output, err := g.CreateOrderUseCase.Execute(input)
					if err != nil {
						return nil, err
					}

					return map[string]interface{}{
						"id":          output.ID,
						"price":       output.Price,
						"tax":         output.Tax,
						"final_price": output.FinalPrice,
					}, nil
				},
			},
		},
	})
}

func (g *Graph) Schema() (graphql.Schema, error) {
	return graphql.NewSchema(graphql.SchemaConfig{
		Query:    g.QueryType(),
		Mutation: g.MutationType(),
	})
}
