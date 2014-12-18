Go Implementation of BioC

Wanli Liu, Rezarta Islamaj Doğan, Dongseop Kwon, Hernani Marques,
Fabio Rinaldi, W. John Wilbur and Donald C. Comeau, "BioC
Implementations in Go, Perl, Python and Ruby", Database, submitted.

United States Government Work
  see PUBLIC DOMAIN NOTICE at end of this file

Install library in Go workspace:
  go get bitbucket.org/comeau/go_bioc/BioC
  go install bitbucket.org/comeau/go_bioc/BioC

Sample programs and BioC files available in tar file
  bioc.sourceforge.net
    go_bioc_1.0.tar.gz

Buid and test sample programs:
  ./build.sh
  ./test.sh    
    
The example programs are not organized in a Go workspace, because they
would require a separate directory for each program. This organization
may change as Go workspaces are understood better.


List of files and what they demonstrate

LICENSE.txt         United States Government Work details
README.txt          this file
bioc_files/         sample BioC files
build.sh            build sample programs
out/                output from sample programs
src/                source of sample programs
test.sh             run sample programs and compare output with
                    expected output

bioc_files:         

BioC.dtd            BioC DTD
collection.xml      typical collection with 10 PubMed documents
everything-sentence.xml       example of every BioC feature with sentences
everything.xml      example of every BioC feature

out:                output from sample programs
collection.out
everything-sentence.out
everything.out

src:                source of sample programs

BioC_Copy.go        copy a BioC file using collection IO
BioC_Copy_Deep.go   copy a BioC file explicilty copying the nested
                    data structures    
BioC_Copy_doc.go    copy a BioC file using document at a time IO
demo.go             demo program from paper showing a cartoon of
                    sample processing 
print.go            print contents of a BioC file
print_serial.go     print contents of a BioC file using document at a
                    time IO

--------------------------------------------------------------------------------
********************************************************************************
--------------------------------------------------------------------------------
|*|  PUBLIC DOMAIN NOTICE
|*|
|*| This work is a "United States Government Work" under the terms of the
|*| United States Copyright Act. It was written as part of the authors'
|*| official duties as a United States Government employee and thus cannot
|*| be copyrighted within the United States. The data is freely available to
|*| the public for use. The National Library of Medicine and the U.S.
|*| Government have not placed any  restriction on its use or reproduction
|*|
|*| Although all reasonable efforts have been taken to ensure the accuracy and
|*| reliability of the data and its source code, the NLM and the U.S. Government
|*| do not and cannot warrant the performance or results that may be obtained by
|*| using it. The NLM and the U.S. Government disclaim all warranties, express
|*| or implied, including warranties of performance, merchantability or fitness
|*| for any particular purpose.
|*|
|*| Please cite the authors in any work or product based on this material:
|*| Wanli Liu, Rezarta Islamaj Doğan, Dongseop Kwon, Hernani Marques,
|*| Fabio Rinaldi, W. John Wilbur and Donald C. Comeau, "BioC
|*| Implementations in Go, Perl, Python and Ruby", Database,
|*| submitted. 
--------------------------------------------------------------------------------
********************************************************************************
--------------------------------------------------------------------------------
