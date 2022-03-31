package main

import (
	"Atomic/handler"
	"fmt"
	"github.com/bengadbois/flippytext"
)

func Menu() {
	fmt.Println(`
           ,ggg,                                                          
          dP""8I     I8                                                   
         dP   88     I8                                                   
        dP    88  88888888                                  gg            
       ,8'    88     I8                                     ""            
       d88888888     I8      ,ggggg,     ,ggg,,ggg,,ggg,    gg     ,gggg, 
 __   ,8"     88     I8     dP"  "Y8ggg ,8" "8P" "8P" "8,   88    dP"  "Yb
dP"  ,8P      Y8    ,I8,   i8'    ,8I   I8   8I   8I   8I   88   i8'      
Yb,_,dP       '8b, ,d88b, ,d8,   ,d8'  ,dP   8I   8I   Yb,_,88,_,d8,_    _
"Y8P"         'Y888P""Y88P"Y8888P"    8P'   8I   8I   'Y88P""Y8P""Y8888PP
	`)
	flippytext.New().Write(`			Made by leki#6796`)

}

func main() {

	Menu()
}
