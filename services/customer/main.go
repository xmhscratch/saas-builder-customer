package main

import (
	"log"

	"localdomain/customer/core"
	corehttp "localdomain/customer/core/http"
)

func main() {
	defer func() {
		if rec := recover(); rec != nil {
			log.Println(rec)
			return
		}
	}()

	cfg, err := core.NewConfig("")
	if err != nil {
		panic(err)
	}

	srv, err := corehttp.NewServer(cfg)
	if err != nil {
		panic(err)
	}

	if err := srv.Start(); err != nil {
		panic(err)
	}
}

// TRUNCATE TABLE `customers`;
// TRUNCATE TABLE `customer_attribute_blob`;
// TRUNCATE TABLE `customer_attribute_boolean`;
// TRUNCATE TABLE `customer_attribute_datetime`;
// TRUNCATE TABLE `customer_attribute_decimal`;
// TRUNCATE TABLE `customer_attribute_int`;
// TRUNCATE TABLE `customer_attribute_text`;
// TRUNCATE TABLE `customer_attribute_varchar`;
