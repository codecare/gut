gut is a command line tool to analyse and improve git repositories

this is still a very early working stage

sample output
        
      analyis of branch XYZ-1705-secure-headers in respect to develop 
      analysis result: FullyMerged XYZ-1705-secure-headers in develop
      merge commit: a82003e7933f5d811abd8475d0efc8e491e1b115: parents [999847a6b2667e9c2deec7d7a230bf87dcdf8a96 e315b5c36c96f4df2c2b0b31b04aed46719739e0], author: stehling <us@codecare.de> 1580225870 +0100, committer:stehling <us@codecare.de> 1580225870 +0100, msg: [Merge branch 'develop' into XYZ-2466-secureheaders]
      top commit of branch: e315b5c36c96f4df2c2b0b31b04aed46719739e0: parents [f9fbf6698249822b4bad1f04031ff5724c9a0ae3], author: stehling <us@codecare.de> 1580222282 +0100, committer:stehling <us@codecare.de> 1580222282 +0100, msg: [XYZ-1705 reporting]
         
use at your own risk

Apache v2 License


you can install it with go get -v github.com/codecare/gut/cmd/gut

this will install an executable in your default path for go executables $GOPATH/bin (default GOPATH=$HOME/go)
