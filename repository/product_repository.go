package repository

import (
	"database/sql"
	"fmt"
	"products-api/model"
)

type ProductRepository struct {
	connection *sql.DB
}

func NewProductRepository(connection *sql.DB) ProductRepository {
	return ProductRepository{
		connection: connection,
	}
}

func (pr *ProductRepository) GetProducts() ([]model.Product, error) {
	query := "SELECT id, product_name, price FROM product"
	rows, err := pr.connection.Query(query)
	if err != nil {
		fmt.Println(err)
		return []model.Product{}, err
	}

	var productList []model.Product
	var productObj model.Product

	for rows.Next() {
		err = rows.Scan(
			&productObj.ID,
			&productObj.Name,
			&productObj.Price)

		if err != nil {
			fmt.Println(err)
			return []model.Product{}, err
		}

		productList = append(productList, productObj)
	}

	rows.Close()

	return productList, nil
}

func (pr *ProductRepository) CreateProduct(product model.Product) (int, error) {
	var id int
	query, err := pr.connection.Prepare("INSERT INTO product" +
		"(product_name, price)" +
		" VALUES ($1, $2) RETURNING id")
	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	err = query.QueryRow(product.Name, product.Price).Scan(&id)

	if err != nil {
		fmt.Println(err)
		return 0, err
	}

	query.Close()

	return id, nil
}

func (pr *ProductRepository) GetProduct(id int) (*model.Product, error) {
	queryString := "SELECT * FROM product WHERE product.id = $1"

	query, err := pr.connection.Prepare(queryString)

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	var product model.Product

	err = query.QueryRow(id).Scan(
		&product.ID,
		&product.Name,
		&product.Price,
	)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		fmt.Println(err)
		return nil, err
	}

	query.Close()
	return &product, nil
}

func (pr *ProductRepository) UpdateProduct(id int, product *model.Product) (*model.Product, error) {
	query, err := pr.connection.Prepare("UPDATE product" +
		" SET product_name = $2, price = $3" +
		" WHERE id = $1" +
		" returning id")

	if err != nil {
		fmt.Println(err)
		return nil, err
	}

	err = query.QueryRow(id, product.Name, product.Price).Scan(&product.ID)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		fmt.Println(err)
		return nil, err
	}

	query.Close()
	return product, nil
}
