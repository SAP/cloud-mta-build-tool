## Troubleshooting

### Timeouts

 - registry - Sometimes, when executing the `mbt build [args]` command, the build process hangs and ends with timeout error. 
   This may be due to a network error or because of missing registry configurations. 
   When using packages provided by SAP, you should verify that your `npm config` file refers to the SAP registry. 
   
   Proposed solution: 
   
   Add an `.npmrc` file to the module (that hangs) as a sibling to the 'package.json' file as follows:
  
```
  $ cat .npmrc

  @sap:registry=https://npm.sap.com/

```

 For more detail's, refer to npm [docs](https://docs.npmjs.com/files/npmrc)
  

   




