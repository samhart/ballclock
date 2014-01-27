package main

import (
		"os"
		"fmt"
	        "bufio"
		"strconv"
		"sync"
		)



//simulates a full 12 hour ball clock cycle
func permute(inqueue *[]int){

	//make a local copy of inqueue
	queue := make([]int, len(*inqueue))
	queue = *inqueue;
	//slices to represent the one min, five min and hour components
	ones := make([]int, 0)
	fives := make([]int, 0)
	hours := make([]int, 0)

	for{
		//pop the first element off the queue
		next := queue[0]
		qt := make([]int, len(queue) - 1)
		copy(qt, queue[1:])
		queue = qt

		//add minutes
		if(len(ones) >= 4){//minute section is full
			//empty the ones into the queue, and add next to the fives
			for oi := len(ones) - 1; oi >= 0; oi--{
				queue = append(queue, ones[oi])
			}
			ones = make([]int,0)
			//empty the fives into the queue and add next to the hours
			if(len(fives) >= 11){//fives is full
				for fi := len(fives) - 1; fi >= 0; fi--{
					queue = append(queue, fives[fi])
				}
				fives = make([]int,0)

				if(len(hours) >= 11){//hours is full, this is the last cycle
					//empty the hours into the queue
					for hi := len(hours) - 1; hi >= 0; hi--{
						queue = append(queue, hours[hi])
					}
					//the added ball gets dumped back into the queue
					queue = append(queue, next)
					hours = make([]int,0)//empty hours
					//update inqueue
					*inqueue = queue
					break
				}else{
					hours = append(hours, next)
				}

			}else{//add next to the fives
				fives = append(fives, next)
			}
		}else{//add next to the ones
			ones = append(ones, next)
		}

	}

}

//helper function to see if slice is in original order
func inOrder(queue []int) bool{
	res := true;
	if(len(queue) < 2){
		return true
	}
	lastNum := 0
	for i := 0; i < len(queue); i++{
		if(lastNum == 0){
			lastNum = queue[i]
			continue
		}
		if(queue[i] - lastNum != 1){
			return false
		}
		lastNum = queue[i]
	}
	return res
}

func clockWorker(jobs <-chan int, wg *sync.WaitGroup){

	for els := range jobs {

		queue := make([]int, 0)

		for ei :=1 ; ei <= els; ei++{
			queue = append(queue, ei)
		}

		clockCycles := 0
		for{
			permute(&queue)
			clockCycles++
			if(inOrder(queue)){
				break
			}

		}
		if(clockCycles % 2 != 0){//round up to even #
			clockCycles++
		}
		 fmt.Println(fmt.Sprintf("%d balls cycle after %d days.",els,clockCycles/2))

	}
	wg.Done()

}



func main() {
	//fileName := "/Users/samhart/go/src/input.txt"
	fileName := os.Args[1]
	//set up the workers
	jobs := make(chan int, 100)


	numWorkers := 8
	var wg sync.WaitGroup
	wg.Add(numWorkers)
	for w := 0; w < numWorkers; w++ {
		go clockWorker(jobs, &wg)
	}


	//open input file
	file, err := os.Open(fileName)
	if err != nil {
		panic(err)
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		j, e := strconv.Atoi(scanner.Text())
		if e != nil {
			panic(err)
		}
		if(j == 0){//end signal, close jobs channel
			close(jobs)
		}else{
			jobs <- j
		}

	}
	//wait for all the goroutines to finish
	wg.Wait()


}
