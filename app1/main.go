package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
)

// working with csv file.
// Create func to read in csv file
// Then print two csv files, with same name 1_ & 2_ appended
// Ex. mycsv.csv  -> 1_mycsv.csv, 2_mycsv.csv
// 1_ must have the name of purchased items, as well as the ratio.
// So if 5 items purchased and 2 of them were "wheels", ratio would be
// 2 / 5, etc for each item.

// 2_ must have the item(s) that were the most popular (total purchases)
// name of item, total units purchased altogether

// CSV structure id, location, item, quantity

type Order struct {
	Item     string
	Quantity int
}

type Orders map[string][]int

func main() {
	filename := os.Args[1]

	file, err := os.Open(filename)
	if err != nil {
		println(err)
	}
	defer file.Close()

	data, err := csv.NewReader(file).ReadAll()
	if err != nil {
		println(err)
	}

	orders := Orders{}
	order_pop := make(map[string]int)
	var ordersTotal int
	checkMax := [][]string{
		{"0", "0"},
	}

	for _, val := range data {
		quant, _ := strconv.Atoi(val[3])
		order := Order{
			Item:     val[2],
			Quantity: quant,
		}
		// Add to Orders map
		orders[order.Item] = append(orders[order.Item], order.Quantity)
		// Add to order_pop map
		order_pop[order.Item] += order.Quantity
		fmt.Println(order_pop[order.Item])
		// Using the checkMax string array, we look at each item once
		// Otherwise you would have to reiterate for this information.
		val, err := strconv.Atoi(checkMax[0][1])
		if err != nil {
			println(err)
		}
		switch {
		case val < order_pop[order.Item]:
			{
				checkMax[0][0] = order.Item
				checkMax[0][1] = strconv.Itoa(order.Quantity)
				fmt.Println("New top val", val, checkMax[0][0], checkMax[0][1])
				val = order_pop[order.Item]

			}
		case val == order_pop[order.Item]:
			{
				t := fmt.Sprintf("%d", order.Quantity)
				tempSlice := []string{order.Item, t}
				checkMax = append(checkMax, tempSlice)
				fmt.Println("Pushed another topper", order.Item, t)
			}
		default:
		}

		ordersTotal++
	}

	base := filepath.Base(filename) // get base path for filename
	err = orders.Write_Ratio(base, ordersTotal-1)
	if err != nil {
		println(err)
	}

	err = Write_Popular(checkMax, base)

}

func (o *Orders) Write_Ratio(base string, orderTotal int) error {
	filename := "1_" + base
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()
	for key, value := range *o {
		oT := float64(orderTotal)
		v := float64(len(value))
		ratio := fmt.Sprintf("%.2f", v/oT*100)
		fileInfo := []string{key, ratio}
		err = writer.Write(fileInfo)
		if err != nil {
			return err
		}
	}
	return nil
}

func Write_Popular(order [][]string, base string) error {
	filename := "2_" + base
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()
	for _, value := range order {
		err = writer.Write(value)
		if err != nil {
			return err
		}
	}
	return nil
}
