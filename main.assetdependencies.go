package main

import (
	"encoding/xml"
	"errors"
	"fmt"

	"github.com/hornbill/pb"
)

//cacheAssetDependencies  - caches asset dependency records from instance
func cacheAssetDependencies() error {
	//Get Count
	var err error
	assetDependencyCount, err := getAssetDependencyCount()
	if err != nil {
		return err
	}

	if assetDependencyCount == 0 {
		logger(1, "No existing asset dependencies could be found", true, true)
		return nil
	}
	var i int
	logger(1, "Retrieving "+fmt.Sprint(assetDependencyCount)+" asset dependencies from Hornbill. Please wait...", true, true)

	bar := pb.New(assetDependencyCount)
	bar.ShowPercent = false
	bar.ShowCounters = false
	bar.ShowTimeLeft = false
	bar.Start()
	for i = 0; i <= assetDependencyCount; i += xmlmcPageSize {
		blockAssetDeps, err := getAssetDependencies(i, xmlmcPageSize)
		if err != nil {
			bar.Finish()
			return err
		}
		if len(blockAssetDeps) > 0 {
			for _, v := range blockAssetDeps {
				concatedAssets := v.LID + ":" + v.RID
				assetDependencies[concatedAssets] = v
			}
		}
		bar.Add(xmlmcPageSize)
	}
	bar.Finish()
	logger(1, fmt.Sprint(len(assetDependencies))+" asset dependencies cached.", true, true)
	return err
}

func getAssetDependencyCount() (int, error) {
	espXmlmc.SetParam("application", "com.hornbill.servicemanager")
	espXmlmc.SetParam("table", "h_cmdb_config_items_dependency")
	espXmlmc.SetParam("where", "h_entity_l_name = 'asset' AND h_entity_r_name = 'asset'")
	if configDryrun {
		logger(3, "[DRYRUN] [DEPENDENCY] [COUNT] "+espXmlmc.GetParam(), false, false)
	}
	xmlAssetLinksCount, err := espXmlmc.Invoke("data", "getRecordCount")
	if err != nil {
		retError := "getAssetDependencyCount:Invoke:" + err.Error()
		return 0, errors.New(retError)
	}

	var xmlResponse methodCallResult
	err = xml.Unmarshal([]byte(xmlAssetLinksCount), &xmlResponse)
	if err != nil {
		retError := "getAssetDependencyCount:Unmarshal:" + err.Error()
		return 0, errors.New(retError)
	}
	if xmlResponse.Status != "ok" {
		retError := "getAssetDependencyCount:Xmlmc:" + xmlResponse.State.ErrorRet
		return 0, errors.New(retError)
	}
	return xmlResponse.Params.Count, err
}

func getAssetDependencies(rowStart, limit int) ([]assetDependencyStruct, error) {
	var assetDependenciesBlock []assetDependencyStruct
	espXmlmc.SetParam("application", "com.hornbill.servicemanager")
	espXmlmc.SetParam("queryName", "getDependencies")
	espXmlmc.OpenElement("queryParams")
	espXmlmc.SetParam("rowstart", fmt.Sprint(rowStart))
	espXmlmc.SetParam("limit", fmt.Sprint(limit))
	espXmlmc.CloseElement("queryParams")
	if configDryrun {
		logger(3, "[DRYRUN] [DEPENDENCY] [GET] "+espXmlmc.GetParam(), false, false)
	}
	xmlAssets, err := espXmlmc.Invoke("data", "queryExec")
	if err != nil {
		retError := "getAssetDependencies:Invoke:" + err.Error()
		return assetDependenciesBlock, errors.New(retError)
	}

	var xmlResponse methodCallResultDependencies
	err = xml.Unmarshal([]byte(xmlAssets), &xmlResponse)
	if err != nil {
		retError := "getAssetDependencies:Unmarshal:" + err.Error()
		return assetDependenciesBlock, errors.New(retError)
	}
	if xmlResponse.Status != "ok" {
		retError := "getAssetDependencies:Xmlmc:" + xmlResponse.State.ErrorRet
		return assetDependenciesBlock, errors.New(retError)
	}
	return xmlResponse.Dependencies, err
}

func addDependency(lid, rid, dependency string) error {
	espXmlmc.SetParam("application", "com.hornbill.servicemanager")
	espXmlmc.SetParam("entity", "ConfigurationItemsDependency")
	espXmlmc.OpenElement("primaryEntityData")
	espXmlmc.OpenElement("record")
	espXmlmc.SetParam("h_entity_l_id", lid)
	espXmlmc.SetParam("h_entity_l_name", "asset")
	espXmlmc.SetParam("h_entity_r_id", rid)
	espXmlmc.SetParam("h_entity_r_name", "asset")
	espXmlmc.SetParam("h_dependency", dependency)
	espXmlmc.CloseElement("record")
	espXmlmc.CloseElement("primaryEntityData")
	if configDryrun {
		logger(3, "[DRYRUN] [DEPENDENCY] [CREATE] "+espXmlmc.GetParam(), false, false)
		espXmlmc.ClearParam()
		return nil
	}
	linkAssetResult, err := espXmlmc.Invoke("data", "entityAddRecord")
	if err != nil {
		retError := "addDependency:Invoke:" + err.Error()
		return errors.New(retError)
	}

	var xmlResponse methodCallResult
	err = xml.Unmarshal([]byte(linkAssetResult), &xmlResponse)
	if err != nil {
		retError := "addDependency:Unmarshal:" + err.Error()
		return errors.New(retError)
	}
	if xmlResponse.Status != "ok" {
		retError := "addDependency:Xmlmc:" + xmlResponse.State.ErrorRet
		return errors.New(retError)
	}
	return nil
}

func updateDependency(id, dependency string) error {
	espXmlmc.SetParam("application", "com.hornbill.servicemanager")
	espXmlmc.SetParam("entity", "ConfigurationItemsDependency")
	espXmlmc.OpenElement("primaryEntityData")
	espXmlmc.OpenElement("record")
	espXmlmc.SetParam("h_pk_confitemdependencyid", id)
	espXmlmc.SetParam("h_dependency", dependency)
	espXmlmc.CloseElement("record")
	espXmlmc.CloseElement("primaryEntityData")
	if configDryrun {
		logger(3, "[DRYRUN] [DEPENDENCY] [UPDATE] "+espXmlmc.GetParam(), false, false)
		espXmlmc.ClearParam()
		return nil
	}
	linkAssetResult, err := espXmlmc.Invoke("data", "entityUpdateRecord")
	if err != nil {
		retError := "updateDependency:Invoke:" + err.Error()
		return errors.New(retError)
	}

	var xmlResponse methodCallResult
	err = xml.Unmarshal([]byte(linkAssetResult), &xmlResponse)
	if err != nil {
		retError := "updateDependency:Unmarshal:" + err.Error()
		return errors.New(retError)
	}
	if xmlResponse.Status != "ok" {
		retError := "updateDependency:Xmlmc:" + xmlResponse.State.ErrorRet
		return errors.New(retError)
	}
	return nil
}

func deleteDependency(id string) error {
	espXmlmc.SetParam("application", "com.hornbill.servicemanager")
	espXmlmc.SetParam("entity", "ConfigurationItemsDependency")
	espXmlmc.SetParam("keyValue", id)
	if configDryrun {
		logger(3, "[DRYRUN] [DEPENDENCY] [DELETE] "+espXmlmc.GetParam(), false, false)
		espXmlmc.ClearParam()
		return nil
	}
	linkAssetResult, err := espXmlmc.Invoke("data", "entityDeleteRecord")
	if err != nil {
		retError := "deleteDependency:Invoke:" + err.Error()
		return errors.New(retError)
	}

	var xmlResponse methodCallResult
	err = xml.Unmarshal([]byte(linkAssetResult), &xmlResponse)
	if err != nil {
		retError := "deleteDependency:Unmarshal:" + err.Error()
		return errors.New(retError)
	}
	if xmlResponse.Status != "ok" {
		retError := "deleteDependency:Xmlmc:" + xmlResponse.State.ErrorRet
		return errors.New(retError)
	}
	return nil
}
