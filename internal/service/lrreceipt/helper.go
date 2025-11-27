package lrreceipt

import (
	"strconv"

	"github.com/patrickmn/go-cache"

	"go-transport-hub/dtos"
)

func (lr *LRReceipt) BuildCityMaps(podLoadingPoints *[]dtos.PodLoadUnLoad, cacheMap *cache.Cache) (map[int64]string, map[int64]string) {

	loadingMap := make(map[int64]string)
	unloadingMap := make(map[int64]string)

	for i := range *podLoadingPoints {
		lp := &(*podLoadingPoints)[i]
		cacheKey := strconv.FormatInt(lp.LoadingPointId, 10)

		var cityCode string

		if v, found := cacheMap.Get(cacheKey); found {
			loc := v.(*dtos.LoadUnLoadLoc)
			lp.CityCode, lp.CityName = loc.CityCode, loc.CityName
			//lr.l.Info("from cache:", loc.CityCode, loc.CityName)
			cityCode = loc.CityCode
		} else {
			loc, err := lr.lrReceiptDao.GetLocationNameById(lp.LoadingPointId)
			if err != nil {
				lr.l.Error("ERROR: loadingPoints", err)
				continue
			}
			lp.CityCode, lp.CityName = loc.CityCode, loc.CityName
			cacheMap.Set(cacheKey, loc, cache.DefaultExpiration)
			cityCode = loc.CityCode
		}

		switch lp.Type {
		case "loading_point":
			loadingMap[lp.TripSheetId] = appendCityCode(loadingMap[lp.TripSheetId], cityCode)
		case "un_loading_point":
			unloadingMap[lp.TripSheetId] = appendCityCode(unloadingMap[lp.TripSheetId], cityCode)
		}
	}

	return loadingMap, unloadingMap
}

func appendCityCode(existing, city string) string {
	if existing == "" {
		return city
	}
	return existing + "-" + city
}
