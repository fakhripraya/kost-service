package handlers

import (
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/fakhripraya/kost-service/config"
	"github.com/fakhripraya/kost-service/data"
	"github.com/fakhripraya/kost-service/database"
	"github.com/fakhripraya/kost-service/entities"
)

// GetKost is a method to fetch the given kost info
func (kostHandler *KostHandler) GetKost(rw http.ResponseWriter, r *http.Request) {

	// get the kost via context
	kostReq := r.Context().Value(KeyKost{}).(*entities.Kost)

	// look for the selected kost in the db to fetch all the picts
	var selectedKost database.DBKost
	if err := config.DB.Where("id = ?", kostReq.ID).First(&selectedKost).Error; err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)

		return
	}

	// parse the given instance to the response writer
	err := data.ToJSON(selectedKost, rw)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)

		return
	}

	return
}

// GetKostPicts is a method to fetch the given kost picts list
func (kostHandler *KostHandler) GetKostPicts(rw http.ResponseWriter, r *http.Request) {

	// get the kost via context
	kostReq := r.Context().Value(KeyKost{}).(*entities.Kost)

	// look for the selected kost in the db to fetch all the picts
	var kostPicts []database.DBKostPict
	if err := config.DB.Where("kost_id = ?", kostReq.ID).Find(&kostPicts).Error; err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)

		return
	}

	// parse the given instance to the response writer
	err := data.ToJSON(kostPicts, rw)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)

		return
	}

	return
}

// GetKostFacilities is a method to fetch the given kost facilities list
func (kostHandler *KostHandler) GetKostFacilities(rw http.ResponseWriter, r *http.Request) {

	// get the kost via context
	kostReq := r.Context().Value(KeyKost{}).(*entities.Kost)

	// look for the selected kost in the db to fetch all the picts
	var kostfacilities []entities.KostFacilities
	if err := config.DB.
		Model(&database.DBKostFacilities{}).
		Select("db_kost_facilities.id, db_kost_facilities.fac_id, db_kost_facilities.kost_id, master_facilities.fac_name as fac_desc").
		Joins("inner join master_facilities on master_facilities.id = db_kost_facilities.fac_id").
		Where("db_kost_facilities.kost_id = ? AND master_facilities.fac_category = ?", kostReq.ID, 0).Scan(&kostfacilities).Error; err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)

		return
	}

	// parse the given instance to the response writer
	err := data.ToJSON(kostfacilities, rw)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)

		return
	}

	return
}

// GetKostList is a method to fetch the given kost info
func (kostHandler *KostHandler) GetKostList(rw http.ResponseWriter, r *http.Request) {

	rw.WriteHeader(http.StatusOK)
	return
}

// GetMyKost is a method to fetch the given kost info
func (kostHandler *KostHandler) GetMyKost(rw http.ResponseWriter, r *http.Request) {

	// get the current user login
	var currentUser *database.MasterUser
	currentUser, err := kostHandler.kost.GetCurrentUser(rw, r, kostHandler.store)
	if err != nil {
		data.ToJSON(&GenericError{Message: err.Error()}, rw)

		return
	}

	// look for the current kost in the db
	var myKost database.DBKost
	if err := config.DB.Where("owner_id = ?", currentUser.ID).First(&myKost).Error; err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)

		return
	}

	// parse the given instance to the response writer
	err = data.ToJSON(myKost, rw)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)

		return
	}

	return
}

// GetMyKostList is a method to fetch the list of the given kost info
func (kostHandler *KostHandler) GetMyKostList(rw http.ResponseWriter, r *http.Request) {

	// get the current user login
	var currentUser *database.MasterUser
	currentUser, err := kostHandler.kost.GetCurrentUser(rw, r, kostHandler.store)
	if err != nil {
		data.ToJSON(&GenericError{Message: err.Error()}, rw)

		return
	}

	// look for the current kost list in the db
	var kostList []database.DBKost
	if err := config.DB.Where("owner_id = ?", currentUser.ID).Find(&kostList).Error; err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)

		return
	}

	// parse the given instance to the response writer
	err = data.ToJSON(kostList, rw)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)

		return
	}

	return
}

// GetEventList is a method to fetch the list of application event
func (kostHandler *KostHandler) GetEventList(rw http.ResponseWriter, r *http.Request) {

	// look for the current kost list in the db
	var eventList []database.MasterEvent
	if err := config.DB.Where("is_active = ?", true).Find(&eventList).Error; err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)

		return
	}

	// parse the given instance to the response writer
	err := data.ToJSON(eventList, rw)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)

		return
	}

	return
}

// GetNearYouList is a method to fetch the list of nearby kost
func (kostHandler *KostHandler) GetNearYouList(rw http.ResponseWriter, r *http.Request) {

	// get the current logged in user additional info via context
	userReq := r.Context().Value(KeyUser{}).(*entities.User)

	baseURL, _ := url.Parse("http://api.positionstack.com")

	baseURL.Path += "v1/reverse"

	params := url.Values{}

	// Access Key
	params.Add("access_key", os.Getenv("GEOCODER_API_KEY"))

	// Query = latitude,longitude
	params.Add("query", userReq.Latitude+","+userReq.Longitude)

	// trigger the reverse geocoder request to fetch the addresses data
	baseURL.RawQuery = params.Encode()
	req, _ := http.NewRequest("GET", baseURL.String(), nil)
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)

		return
	}

	defer res.Body.Close()

	// create the geo location instance
	geoLocation := &entities.Geolocation{}
	data.FromJSON(geoLocation, res.Body)

	// look for the current kost list in the db
	var nearbyKostList []database.DBKost
	if err := config.DB.Limit(20).Where("is_active = ? AND city >= ?", true, geoLocation.GeoData[0].County).Find(&nearbyKostList).Error; err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)

		return
	}

	type Ranges struct {
		KostID   uint
		Distance float64
	}
	type NearbyKostView struct {
		ID             uint   `json:"id"`
		KostName       string `json:"kost_name"`
		City           string `json:"city"`
		ThumbnailURL   string `json:"thumbnail_url"`
		ThumbnailPrice string `json:"thumbnail_price"`
	}
	type FinalNearbyKostView struct {
		CarouselList []NearbyKostView `json:"carousel_list"`
	}

	var listRanges []Ranges
	var tempRanges []Ranges

	for _, nearby := range nearbyKostList {

		lat1, _ := strconv.ParseFloat(userReq.Latitude, 64)
		lng1, _ := strconv.ParseFloat(userReq.Longitude, 64)
		lat2, _ := strconv.ParseFloat(nearby.Latitude, 64)
		lng2, _ := strconv.ParseFloat(nearby.Longitude, 64)

		distance := kostHandler.kost.CalculateDistanceBetween(lat1, lng1, lat2, lng2, "K")
		tempRanges = append(tempRanges, Ranges{
			KostID:   nearby.ID,
			Distance: distance,
		})

	}

	for i := 0; i < 5; i++ {

		var smallest float64 = -1
		for _, obj := range tempRanges {

			if smallest == -1 {
				smallest = obj.Distance
			} else {
				if smallest > obj.Distance {
					smallest = obj.Distance
				}
			}
		}

		var keepRanges []Ranges

		for _, num := range tempRanges {

			if num.Distance == smallest {
				listRanges = append(listRanges, num)
			} else {
				keepRanges = append(keepRanges, num)
			}
		}

		tempRanges = keepRanges
	}

	// variable to hold temporary data of nearby kost list
	var tempNearbyKostList []NearbyKostView
	var finalNearbyKostList []FinalNearbyKostView
	var i = 0

	for _, nearby := range listRanges {

		var selectedKost database.DBKost
		if err := config.DB.Where("id= ?", nearby.KostID).First(&selectedKost).Error; err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			data.ToJSON(&GenericError{Message: err.Error()}, rw)

			return
		}

		var roomPrice struct {
			RoomPrice    float64 `json:"room_price"`
			RoomPriceUom uint    `json:"room_price_uom"`
		}

		if err := config.DB.Raw("SELECT room_price, room_price_uom FROM db_kost_rooms WHERE kost_id = ?", nearby.KostID).Scan(&roomPrice).Error; err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			data.ToJSON(&GenericError{Message: err.Error()}, rw)

			return
		}

		var uomDescription string
		if err := config.DB.Raw("SELECT uom_desc FROM master_uoms WHERE id = ?", roomPrice.RoomPriceUom).Scan(&uomDescription).Error; err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			data.ToJSON(&GenericError{Message: err.Error()}, rw)

			return
		}

		if i < 3 {
			tempNearbyKostList = append(tempNearbyKostList, NearbyKostView{
				ID:             nearby.KostID,
				KostName:       selectedKost.KostName,
				City:           selectedKost.City,
				ThumbnailURL:   selectedKost.ThumbnailURL,
				ThumbnailPrice: fmt.Sprintf("%f", roomPrice.RoomPrice) + " / " + uomDescription,
			})

			i++
		} else {
			finalNearbyKostList = append(finalNearbyKostList, FinalNearbyKostView{
				CarouselList: tempNearbyKostList,
			})

			var tempData = tempNearbyKostList[2]

			tempNearbyKostList = nil
			tempNearbyKostList = append(tempNearbyKostList, tempData)
			tempNearbyKostList = append(tempNearbyKostList, NearbyKostView{
				ID:             nearby.KostID,
				KostName:       selectedKost.KostName,
				City:           selectedKost.City,
				ThumbnailURL:   selectedKost.ThumbnailURL,
				ThumbnailPrice: fmt.Sprintf("%f", roomPrice.RoomPrice) + " / " + uomDescription,
			})

			i = 2
		}

	}

	if len(tempNearbyKostList) > 0 {
		finalNearbyKostList = append(finalNearbyKostList, FinalNearbyKostView{
			CarouselList: tempNearbyKostList,
		})
	}

	if len(finalNearbyKostList[len(finalNearbyKostList)-1].CarouselList) == 3 {

		var tempLastNearby []NearbyKostView
		tempLastNearby = append(tempLastNearby, NearbyKostView{
			ID:             finalNearbyKostList[len(finalNearbyKostList)-1].CarouselList[len(finalNearbyKostList[len(finalNearbyKostList)-1].CarouselList)-1].ID,
			KostName:       finalNearbyKostList[len(finalNearbyKostList)-1].CarouselList[len(finalNearbyKostList[len(finalNearbyKostList)-1].CarouselList)-1].KostName,
			City:           finalNearbyKostList[len(finalNearbyKostList)-1].CarouselList[len(finalNearbyKostList[len(finalNearbyKostList)-1].CarouselList)-1].City,
			ThumbnailURL:   finalNearbyKostList[len(finalNearbyKostList)-1].CarouselList[len(finalNearbyKostList[len(finalNearbyKostList)-1].CarouselList)-1].ThumbnailURL,
			ThumbnailPrice: finalNearbyKostList[len(finalNearbyKostList)-1].CarouselList[len(finalNearbyKostList[len(finalNearbyKostList)-1].CarouselList)-1].ThumbnailPrice,
		})
		finalNearbyKostList = append(finalNearbyKostList, FinalNearbyKostView{
			CarouselList: tempLastNearby,
		})
	}

	// parse the given instance to the response writer
	err = data.ToJSON(finalNearbyKostList, rw)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)

		return
	}

	return
}
