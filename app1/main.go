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
// 1_ must have the name of purchased items, as well as the percentage for each.
// So if 5 items purchased and 2 of them were "wheels", percentage would be
// 2 / 5 * 100, for each item.

// 2_ must have the item(s) that were the most popular (total purchases)
// name of item, total units purchased altogether

// CSV structure id, location, item, quantity

type Order struct {
	Item     string
	Quantity int
}

type Orders map[string][]int
type Popular map[string]int

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

	var ordersTotal int
	orders := Orders{}
	order_pop := Popular{}

	for _, val := range data {
		quant, _ := strconv.Atoi(val[3])
		// Really don't need to do this conversion but it is a little easier to read
		order := Order{
			Item:     val[2],
			Quantity: quant,
		}
		// Add to Orders map
		orders[order.Item] = append(orders[order.Item], order.Quantity)
		// Add to order_pop map val
		order_pop[order.Item] += order.Quantity
		ordersTotal++
	}
	base := filepath.Base(filename) // get base path for filename

	// Ideally we put these in goroutine
	err = Write_Ratio(&orders, base, ordersTotal-1)
	if err != nil {
		fmt.Println(err)
	}
	err = Write_Popular(order_pop, filename)
	if err != nil {
		fmt.Println(err)
	}

}

func Write_Ratio(o *Orders, base string, orderTotal int) error {
	fmt.Println("Writing file now")
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

func Write_Popular(o Popular, base string) error {
	var popularSlice []string // Hold the popular keys
	var checkPopular string
	var count int

	for key, val := range o {
		if len(popularSlice) == 0 {
			popularSlice = append(popularSlice, key)
			checkPopular = popularSlice[0]

		}
		// If we have multiple ties, ideally we sort this to keep from adding
		// values that are not really popular over time
		// *** ISSUES WITH THIS CODE CURRENTLY
		if count > 2 {
			var temp int
			for i := 0; i < count-1; i++ {
				if o[popularSlice[i]] < o[popularSlice[i+1]] {
					popularSlice = append(popularSlice[:i], popularSlice[i+1:]...)
					count--
				}
				if o[popularSlice[i]] > o[popularSlice[i+1]] {
					if temp = i + 2; popularSlice[temp] != "" {
						popularSlice = append(popularSlice[:i+1], popularSlice[temp:]...)
						count--
					} else {
						// need to remove the item but
						popularSlice = append(popularSlice[:i+1], popularSlice[temp:]...)
						count--
					}
				}

			}
		}

		if o[checkPopular] == val {
			popularSlice = append(popularSlice, key)
			count++
		}
		if o[checkPopular] < val {
			//popular = append(popular[0:1], key)
			//Overwrite current val
			popularSlice[0] = key
		}

	}
	filename := "2_" + base
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for i := 0; i < len(popularSlice); i++ {
		key := popularSlice[i]
		val := strconv.Itoa(o[key])
		fileinfo := []string{key, val}
		err := writer.Write(fileinfo)
		if err != nil {
			return err
		}
	}
	return nil
}
