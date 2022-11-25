#!/usr/bin/env sh

  # MTA build process steps

  # ----Pre build process ---------

  # Set the current project repository path for general mta process
    DIR=$(pwd)
  # Create temporary dir for build purpose
    tmpDir=$(mit execute prepare)
  # Copy module folder to temp directory
    mit execute copy ui5app $tmpDir ui5app

   # Change to the temporary folder path
    cd $tmpDir/ui5app

  # ----Executing build for module ui5app -------
  # installing module dependencies & execute grunt & remove dev dependencies
    (npm install && grunt && npm prune production ) &
  # wait to the process to finish
    wait
  # Pack module after build for deployment
    mit execute pack $tmpDir ui5app ui5app
 
   # Move to MTA project level
    cd $DIR

  # ----Post build process------

  # Create META-INF folder with MANIFEST.MF & mtad.yaml
    mit execute meta $tmpDir &
  # Pack as MTAR artifact
    wait
    mit execute mtar $tmpDir $DIR
  # Remove tmp folder
    mit execute cleanup $tmpDir

