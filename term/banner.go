package term

import "fmt"

func PrintBanner() {
  banner := `    
    |\__/,|   ('\
  _.|o o  |_   ) )
-(((---(((----- ezcat // ngn
`

  r := 255             
  b := 0               
  
  for i := 0; i < len(banner); i++ {
    g := 165 - i*(165/len(banner))

    color := fmt.Sprintf("\x1b[38;2;%d;%d;%dm", r, g, b)
    fmt.Printf("%s%c", color, banner[i])
	}
}
