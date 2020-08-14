package cmd

import (
	"bytes"
	"log"
	"os/exec"
	"strings"
)
func getMBPUCount() int {
	var err error
  
	// Create an lspci command.
	  lspci := exec.Command("lspci")
  
	  // Create a grep command that searches for anything
	  // that contains Xilinx in it's filename.
	grep := exec.Command("grep", "Xilinx")
	
	// Set grep's stdin to the output of the lspci command.
	  grep.Stdin, err = lspci.StdoutPipe()
	  if err != nil {
	  log.Fatalln(err)
	  return 0
	  }
  
	// Set grep's stdout to cmdOutput
	cmdOutput := &bytes.Buffer{}
	  grep.Stdout = cmdOutput
  
	// Start the grep command first. (The order will be last command first)
	err = grep.Start()
	if err != nil {
	  log.Fatalln(err)
	  return 0
	  }
  
	  // Run the ls command. (Run calls start and also calls wait)
	err = lspci.Run()
	if err != nil {
	  log.Fatalln(err)
	  return 0
	  }
  
	// Wait for the grep command to finish.
	err = grep.Wait()
	if err != nil {
	  // nothing searched... not printing error just return 0
	  return 0
	  }
	  
	count := len(strings.Split(cmdOutput.String(), "\n"))
	
	return count - 1
  }