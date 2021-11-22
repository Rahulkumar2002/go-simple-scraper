package main 

import(
	"fmt" //to console log or print content on website
	"net/http" //to make http requests
	"os" //to give input from console 
	"strings" //to stringify our tokens 
	"golang.org/x/net/html" //to fetch links data from a tag which is scrapped from our website .
)

func getHref(tz html.Token) (ok bool , href string ){
	for  _ , a := range tz.Attr{
		if a.Key == "href"{
			href = a.Val
		  ok = true
		}
	}
	return 
}
func crawl( url string , ch chan string , chFinished chan bool){
	resp , err := http.Get(url) 

	//we use 'defer' keyword if we want a function to run at the end of the outer function . 
 defer func(){
	 chFinished <- true 
 }()

 //If there is any error in crawling of any url we are checking for it here .

   if err != nil {
      fmt.Println("Error : Failed to crawl : " , url )
     return
   }

   //We are storing response which we have gotten from http.Get() into a variable b 
   b:= resp.Body
   // Closing b after completion of crawl function
        defer b.Close()


    // "golang.org/x/net/html" this package allows us to divide a html page into small tokens .
 //NewTokenizer will distribute this html page which we have gotten as a response from our url .
	z := html.NewTokenizer(b) 

	//After converting our html page into small tokens we will go through each one of them for this we are using a for loop without any checks or variable .

	for{
		//Next funciton is used to move to the next token when we have proccesd the current token .
		tt := z.Next()

		//Using switch case to check between ErrorToken , StartTagToken and EndTagToken.

		switch{
		case tt == html.ErrorToken:
			return 

		case tt == html.StartTagToken : 
		t := z.Token()

		isAnchor := t.Data == "a"
		if !isAnchor{
			continue
		}  

		ok , url := getHref(t)

		if !ok{
			continue
		}

		hasProto := strings.Index(url , "http") == 0
		if hasProto{
			ch <- url
		}
		}
	}
  

}

func main(){
  foundUrls := make(map[string]bool) //Createad a map of keys of type string and value of type bool .
  seedUrls := os.Args[1:] // Used os package to get the arguments from console or shell .

  chUrls := make(chan string)//Created a channel of type string and of name chUrls.
  chFinished := make(chan bool)//Created a channel of name chFinished and of type boolean .

  //Ranging over seedUrls using for loop and storing the URLS in url variable .
  for _, url := range seedUrls{  
	  go crawl( url , chUrls , chFinished) //passing url , chUrls , chFinished into crawl go routine . crawl will go thorugh our porvided link and check for all the unique urls and store it in chUrls and will pass true for unique urls in chFinished .
  }

  //Running a for loop , for to get the value of channel into url and if we have found unique url then we will pass true in chFinished .
  for c := 0 ; c < len(seedUrls) ; {
	  select{
	  case url := <-chUrls  : //Passing value in url from chUrls channel .
		foundUrls[url] = true //If we have found a unique url then make this foundUrls map value true .
	  case <-chFinished : //If chFinished have true value stored then do c++ , c is used in for loop .
	       c++ //incrementing c by one , can also be wriiten as c = c + 1 ;
	  }
  }

  fmt.Println("\nFound" , len(foundUrls) , "unique urls: \n") //We are priniting url which we have found through crawl function.
 
  //We are ranging over foundUrls map and printing all stored urls inside it .
  for url , _ := range foundUrls{   
	  fmt.Println("-" + url)
  }

  //We have closed chUrls channel after it's use .
  close(chUrls)
}