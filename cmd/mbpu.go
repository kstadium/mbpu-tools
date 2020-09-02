package cmd

import (
  "fmt"
  "math/big"
  "strconv"

  "github.com/spf13/cobra"
  "github.com/the-medium/mediumpk"
)

var (
  strD = "519b423d715f8b581f4fa8ee59f4771a5b44c8130b4e3eacca54a56dda72b464"
  strK = "94a1bbb14b906a61a280f245f9e93c7f3b4a6247824f5d33b9670787642a68de"
  strH1 = "ea5cd45052849c4ae816bbc44ed833e832af8a619ba47268aabca2744c4c6268"
  strRExpected = "f3ac8061b514795b8843e3d6629527ed2afd6b1f6a555a7acabb5e6f79c8c2ac"
  strSExpected = "6e9a1aee9981cc4a102aa7033fdf633b39be438527865373edfe90f2ea9e29ac"
  strX = "e305d41ab27b39c84230ab2faf34fb15e9d0543f4ac19d2520b94d71df9be5bf"
  strY = "0b97c506c163237d6e9264f7148336e524d32174754198066995a252b1a51f4e"
  strR = "5806c2774086b61c97afd87585215c09fe57233f232278c0e8976d35f0570641"
  strS = "6d8a758eb8edfeecbdab2e413bee8bc73a88a887f97a54c2a967de0afcb8b0af"
  strH2 = "00000000000000000000000000000000000000000048656c6c6f20576f726c64"
)

var (
  cmdMBPU = &cobra.Command{
    Use:   "mbpu [command]",
    Args: cobra.MinimumNArgs(1),
    Short: "MBPU tool ",
  }

  cmdMBPUVersion = &cobra.Command{
    Use:   "version [device index]",
    Args: cobra.MinimumNArgs(1),
    Short: "Print MBPU version",
    Run: func(cmd *cobra.Command, args []string) {
      version(args[0])      
    },
  }

  cmdMBPUTest = &cobra.Command{
    Use:   "test",
    Short: "Test MBPU",
    Run: func(cmd *cobra.Command, args []string) {
      count := getMBPUCount()
      fmt.Printf("Total MBPU Number : %d\n", count)
      for i := 0; i < count; i++{
        fmt.Printf("MBPU-%d check...\n", i)
        testMBPU(i)    
      }
    },
  }
)

func version(index string) {
  idx, err := strconv.Atoi(index)
  if err != nil {
    fmt.Printf("invalid index : %s\n", index)
    return
  }

  mpk, err := mediumpk.New(idx, 64, "")
  if err != nil {
    fmt.Println(err)
    return
  }

  v, err := mpk.GetVersion()
  if err != nil {
    fmt.Println(err)
    return
  }
  fmt.Printf("%s",v)
  mpk.Close()
  return 
}

func testMBPU(index int) {
  var err error
  fmt.Printf("\tecdsa sign function test : ")
  err = sign(index, 64)
  if err != nil {
    fmt.Printf("failed\n")
    fmt.Printf("\t\t err : %s", err.Error())
  }else {
    fmt.Printf("passed\n")
  }
  
  fmt.Printf("\tecdsa verify function test : ")
  err = verify(index, 64)
  if err != nil {
    fmt.Printf("failed\n")
    fmt.Printf("\t\t err : %s", err.Error())
  }else {
    fmt.Printf("passed\n")
  }

  return
}

func sign(index, maxPending int) (err error) {
  mpk, err := mediumpk.New(index, maxPending, "")
  if err != nil{
    return
  }

  // D
	D := new(big.Int); D.SetString(strD, 16)

	// K
	K := new(big.Int); K.SetString(strK, 16)
	
	// hash
	H := new(big.Int); H.SetString(strH1, 16)
	
	// mediumpk sign
	d32 := make([]byte, 32)
	k32 := make([]byte, 32)
	h32 := make([]byte, 32)
	copy(d32[32-len(D.Bytes()):], D.Bytes()[:])
	copy(k32[32-len(K.Bytes()):], K.Bytes()[:])
  copy(h32[32-len(H.Bytes()):], H.Bytes()[:])
  
  chList := [](chan mediumpk.ResponseEnvelop){}

  for i := 0; i < maxPending; i++ {
		channel := make(chan mediumpk.ResponseEnvelop, 1)
		var req mediumpk.RequestEnvelop = mediumpk.SignRequestEnvelop{d32, k32, h32}
		err = mpk.Request(&channel, req)
		if err != nil {
      return
    }
		chList = append(chList, channel)
  }

  for i := 0; i < maxPending; i++{
		err = mpk.GetResponseAndNotify()
		if err != nil {
      return
    }
  }

	var resp mediumpk.ResponseEnvelop
  for _, v := range chList{
    resp = <- v
    if resp.Result() == 1 {
      fmt.Printf("%d\n", resp.Result())
    }
    close(v)
  }
 
  mpk.Close()

  return
}

func verify(index, maxPending int) (err error) {
  mpk, err := mediumpk.New(index, maxPending, "")
  if err != nil{
    return
  }

  // qx, qy
	X := new(big.Int); X.SetString(strX, 16)
	Y := new(big.Int); Y.SetString(strY, 16)

	// hash
	H := new(big.Int); H.SetString(strH2, 16)
	
	// r, s
	R := new(big.Int); R.SetString(strR, 16)
	S := new(big.Int); S.SetString(strS, 16)
	

	// mediumpk verify
	qx32 := make([]byte, 32)
	qy32 := make([]byte, 32)
	r32 := make([]byte, 32)
	s32 := make([]byte, 32)
	h32 := make([]byte, 32)
	copy(qx32[32-len(X.Bytes()):], X.Bytes()[:])
	copy(qy32[32-len(Y.Bytes()):], Y.Bytes()[:])
	copy(r32[32-len(R.Bytes()):], R.Bytes()[:])
	copy(s32[32-len(S.Bytes()):], S.Bytes()[:])
	copy(h32[32-len(H.Bytes()):], H.Bytes()[:])
  
  chList := [](chan mediumpk.ResponseEnvelop){}

  for i := 0; i < maxPending; i++ {
		channel := make(chan mediumpk.ResponseEnvelop, 1)
		var req mediumpk.RequestEnvelop = mediumpk.VerifyRequestEnvelop{qx32, qy32, r32, s32, h32}
		err = mpk.Request(&channel, req)
		if err != nil {
      return
    }
		chList = append(chList, channel)
  }

  for i := 0; i < maxPending; i++{
		err = mpk.GetResponseAndNotify()
		if err != nil {
      return
    }
  }

	var resp mediumpk.ResponseEnvelop
  for _, v := range chList{
    resp = <- v
    if resp.Result() == 1 {
      fmt.Printf("%d\n", resp.Result())
    }
    close(v)
  }
 
  mpk.Close()

  return

}
