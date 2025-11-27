package commonsvc

import (
	"strconv"
	"time"

	"github.com/patrickmn/go-cache"

	"go-transport-hub/dtos"
)

var lodingPointCityCacheMap = cache.New(15*time.Minute, 30*time.Minute)

func (lr *PreRequisiteObj) BuildCityMaps(loadUnLoadPoints *[]dtos.LoadUnLoadObj) (map[int64][]dtos.LoadUnLoadLoc, map[int64][]dtos.LoadUnLoadLoc) {

	loadingMap := make(map[int64][]dtos.LoadUnLoadLoc)
	unloadingMap := make(map[int64][]dtos.LoadUnLoadLoc)

	for i := range *loadUnLoadPoints {
		lp := &(*loadUnLoadPoints)[i]
		cacheKey := strconv.FormatInt(lp.LoadingPointId, 10)

		if v, found := lodingPointCityCacheMap.Get(cacheKey); found {
			loc := v.(*dtos.LoadUnLoadLoc)
			//lr.l.Debug("from cache:", loc.CityCode, loc.CityName)

			switch lp.Type {
			case "loading_point":
				loadingMap[lp.TripSheetId] = append(loadingMap[lp.TripSheetId], *loc)
			case "un_loading_point":
				unloadingMap[lp.TripSheetId] = append(unloadingMap[lp.TripSheetId], *loc)
			}

		} else {
			loc, err := lr.preRequisiteDao.GetLocationNameById(lp.LoadingPointId)
			if err != nil {
				lr.l.Error("ERROR: loadingPoints", err)
				continue
			}
			lp.CityCode, lp.CityName = loc.CityCode, loc.CityName
			lodingPointCityCacheMap.Set(cacheKey, loc, cache.DefaultExpiration)
			switch lp.Type {
			case "loading_point":
				loadingMap[lp.TripSheetId] = append(loadingMap[lp.TripSheetId], *loc)
			case "un_loading_point":
				unloadingMap[lp.TripSheetId] = append(unloadingMap[lp.TripSheetId], *loc)
			}
		}

	}

	return loadingMap, unloadingMap
}
