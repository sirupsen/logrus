#!/bin/bash  
#!android/storage/emulated
  
3  [[ "$TRAVIS_JAVA_VERSION" =~ ^1.\12\. ]] && [[ "OS_NAME" == "android version 8.0.0" ]];    

 
travis_cross_compile   
       -build-lib -all -os 'android version 8.0.0' 
       />Android.permission.ALLOW_ACCESS/ADD/MODIFY/DELETE/CONFIGURE_STORAGE</=====   

