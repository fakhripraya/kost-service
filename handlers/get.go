package handlers

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/fakhripraya/kost-service/config"
	"github.com/fakhripraya/kost-service/data"
	"github.com/fakhripraya/kost-service/database"
	"github.com/fakhripraya/kost-service/entities"
	"github.com/gorilla/mux"
	"gorm.io/gorm"
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

	var err error

	// get the kost via context
	kostReq := r.Context().Value(KeyKost{}).(*entities.Kost)

	// get the additional field
	vars := mux.Vars(r)
	var facTableName string
	var facTableKey string

	var id uint
	var model *gorm.DB
	var kostFacilities []entities.KostFacilities
	var kostRoomFacilities []entities.KostRoomFacilities

	// filter whether id empty or not
	if vars["roomId"] != "" {

		facTableKey = "room_id"
		facTableName = "db_kost_room_facilities"
		model = config.DB.Model(&database.DBKostRoomFacilities{})

		roomID, err := strconv.ParseUint(vars["roomId"], 10, 32)
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			data.ToJSON(&GenericError{Message: "Unable to convert id"}, rw)

			return
		}

		id = uint(roomID)

	} else {

		id = kostReq.ID
		facTableKey = "kost_id"
		facTableName = "db_kost_facilities"
		model = config.DB.Model(&database.DBKostFacilities{})

	}

	// if id is not empty
	// query will execute a select sql statement towards db_kost_room_facilities
	// if not it will go towards db_kost_facilities instead
	finalQuery := model.
		Select(facTableName+".id,"+facTableName+".fac_id,"+facTableName+"."+facTableKey+",master_facilities.fac_category as fac_category, master_facilities.fac_name as fac_desc").
		Joins("inner join master_facilities on master_facilities.id = "+facTableName+".fac_id").
		Where(facTableName+"."+facTableKey+" = ?", id)

	if vars["roomId"] != "" {
		finalQuery = finalQuery.Scan(&kostRoomFacilities)
	} else {
		finalQuery = finalQuery.Scan(&kostFacilities)
	}

	if err := finalQuery.Error; err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)

		return
	}

	// parse the given instance to the response writer
	if vars["roomId"] != "" {
		err = data.ToJSON(kostRoomFacilities, rw)
	} else {
		err = data.ToJSON(kostFacilities, rw)
	}

	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)

		return
	}

	return
}

// GetKostBenchmark is a method to fetch the given kost benchmark list
func (kostHandler *KostHandler) GetKostBenchmark(rw http.ResponseWriter, r *http.Request) {

	// get the kost via context
	kostReq := r.Context().Value(KeyKost{}).(*entities.Kost)

	// look for the selected kost in the db to fetch all the benchmark
	var kostBenchmark []database.DBKostBenchmark
	if err := config.DB.Where("kost_id = ?", kostReq.ID).Find(&kostBenchmark).Error; err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)

		return
	}

	// parse the given instance to the response writer
	err := data.ToJSON(kostBenchmark, rw)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)

		return
	}

	return
}

// GetKostAccessibility is a method to fetch the given kost accessibility list
func (kostHandler *KostHandler) GetKostAccessibility(rw http.ResponseWriter, r *http.Request) {

	// get the kost via context
	kostReq := r.Context().Value(KeyKost{}).(*entities.Kost)

	// look for the selected kost in the db to fetch all the accessibility
	var kostAccess []database.DBKostAccess
	if err := config.DB.Where("kost_id = ?", kostReq.ID).Find(&kostAccess).Error; err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)

		return
	}

	// parse the given instance to the response writer
	err := data.ToJSON(kostAccess, rw)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)

		return
	}

	return
}

// GetKostAround is a method to fetch the given kost around list
func (kostHandler *KostHandler) GetKostAround(rw http.ResponseWriter, r *http.Request) {

	// get the kost via context
	kostReq := r.Context().Value(KeyKost{}).(*entities.Kost)

	// look for the selected kost in the db to fetch all the around landmark
	var kostAround []database.DBKostAround
	if err := config.DB.Where("kost_id = ?", kostReq.ID).Find(&kostAround).Error; err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)

		return
	}

	// parse the given instance to the response writer
	err := data.ToJSON(kostAround, rw)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)

		return
	}

	return
}

// GetKostReviewList is a method to fetch the given kost review list
func (kostHandler *KostHandler) GetKostReviewList(rw http.ResponseWriter, r *http.Request) {

	// get the kost via context
	kostReq := r.Context().Value(KeyKost{}).(*entities.Kost)

	var kostReview []entities.KostReview
	if err := config.DB.
		Model(&database.DBKostReview{}).
		Select("db_kost_reviews.id, db_kost_reviews.cleanliness,db_kost_reviews.convenience,db_kost_reviews.security,db_kost_reviews.facilities,db_kost_reviews.comments,master_users.display_name, master_users.profile_picture").
		Joins("inner join master_users on master_users.id = db_kost_reviews.user_id").
		Where("db_kost_reviews.kost_id = ?", kostReq.ID, 0).Scan(&kostReview).Error; err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)

		return
	}

	// parse the given instance to the response writer
	err := data.ToJSON(kostReview, rw)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)

		return
	}

	return
}

// GetKostOwner is a method to fetch the given kost owner
func (kostHandler *KostHandler) GetKostOwner(rw http.ResponseWriter, r *http.Request) {

	// get the kost via context
	kostReq := r.Context().Value(KeyKost{}).(*entities.Kost)

	// make a custom model
	type KostOwner struct {
		ID             uint   `json:"id"`
		DisplayName    string `json:"display_name"`
		ProfilePicture string `json:"profile_picture"`
		City           string `json:"city"`
	}

	// look for the selected kost in the db to fetch the owner info aswell
	var kostOwner KostOwner
	if err := config.DB.
		Model(&database.DBKost{}).
		Select("db_kosts.owner_id as id"+
			",master_users.display_name"+
			",master_users.profile_picture"+
			",master_users.city").
		Joins("inner join master_users on master_users.id = db_kosts.owner_id").
		Where("db_kosts.id = ?", kostReq.ID).Scan(&kostOwner).Error; err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)

		return
	}

	// look for the selected owner kost list in the db
	kostList, err := kostHandler.kost.GetKostListByOwner(kostOwner.ID)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)

		return
	}

	finalKostAround := struct {
		ID             uint              `json:"id"`
		DisplayName    string            `json:"display_name"`
		ProfilePicture string            `json:"profile_picture"`
		City           string            `json:"city"`
		KostList       []database.DBKost `json:"kost_list"`
	}{
		ID:             kostOwner.ID,
		DisplayName:    kostOwner.DisplayName,
		ProfilePicture: kostOwner.ProfilePicture,
		City:           kostOwner.City,
		KostList:       kostList,
	}

	// parse the given instance to the response writer
	err = data.ToJSON(finalKostAround, rw)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)

		return
	}

	return
}

// GetKostRoomList is a method to fetch the given kost room list
func (kostHandler *KostHandler) GetKostRoomList(rw http.ResponseWriter, r *http.Request) {

	// get the kost via context
	kostReq := r.Context().Value(KeyKost{}).(*entities.Kost)

	var kostRoom []entities.KostRoom

	if err := config.DB.
		Model(&database.DBKostRoom{}).
		Select("db_kost_rooms.id"+
			",db_kost_rooms.kost_id "+
			",db_kost_rooms.room_desc"+
			",db_kost_rooms.room_price"+
			",db_kost_rooms.room_price_uom"+
			",price.uom_desc as room_price_uom_desc"+
			",db_kost_rooms.room_length"+
			",db_kost_rooms.room_width"+
			",db_kost_rooms.room_area"+
			",db_kost_rooms.room_area_uom"+
			",area.uom_desc as room_area_uom_desc"+
			",db_kost_rooms.max_person"+
			",db_kost_rooms.floor_level"+
			",db_kost_rooms.allowed_gender"+
			",db_kost_rooms.comments"+
			",db_kost_rooms.is_active").
		Joins("inner join master_uoms as area on area.id = db_kost_rooms.room_area_uom").
		Joins("inner join master_uoms as price on price.id = db_kost_rooms.room_price_uom").
		Where("db_kost_rooms.kost_id = ? AND (area.uom_type = ? AND price.uom_type = ?)", kostReq.ID, "length", "currency").Scan(&kostRoom).Error; err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)

		return
	}

	// parse the given instance to the response writer
	err := data.ToJSON(kostRoom, rw)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)

		return
	}

	return
}

// GetKostRoomInfo is a method to fetch the given kost room detail list
func (kostHandler *KostHandler) GetKostRoomInfo(rw http.ResponseWriter, r *http.Request) {

	// get the kost via mux
	vars := mux.Vars(r)
	roomID, err := strconv.ParseUint(vars["roomId"], 10, 32)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		data.ToJSON(&GenericError{Message: "Unable to convert id"}, rw)

		return
	}

	kostRoomDetails, err := kostHandler.kost.GetKostRoomDetails(uint(roomID))
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)

		return
	}

	kostRoomPicts, err := kostHandler.kost.GetKostRoomPicts(uint(roomID))
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)

		return
	}

	kostRoomBookedList, err := kostHandler.kost.GetKostRoomBookedList(uint(roomID))
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)

		return
	}

	var finalKostRoomBookedList []database.DBTransactionRoomBook
	for _, validBooked := range kostRoomBookedList {

		if validBooked.Status == 2 {
			finalKostRoomBookedList = append(finalKostRoomBookedList, validBooked)
		}

	}

	kostDetailView := struct {
		RoomPicts   []database.DBKostRoomPict        `json:"room_picts"`
		RoomDetails []database.DBKostRoomDetail      `json:"room_details"`
		RoomBooked  []database.DBTransactionRoomBook `json:"room_booked"`
	}{
		RoomPicts:   kostRoomPicts,
		RoomDetails: kostRoomDetails,
		RoomBooked:  finalKostRoomBookedList,
	}

	// parse the given instance to the response writer
	err = data.ToJSON(kostDetailView, rw)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)

		return
	}

	return
}

// GetKostList is a method to fetch the given kost info
func (kostHandler *KostHandler) GetKostList(rw http.ResponseWriter, r *http.Request) {

	// get the page via mux
	vars := mux.Vars(r)
	page, err := strconv.Atoi(vars["page"])
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		data.ToJSON(&GenericError{Message: "Unable to convert value"}, rw)

		return
	}

	// look for the current kost list in the db
	kostList, err := kostHandler.kost.GetKostList(page)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)

		return
	}

	var finalKostList []database.DBKost
	for index, kost := range kostList {

		if index >= ((page * 10) - 10) {
			finalKostList = append(finalKostList, kost)
		}

	}

	// parse the given instance to the response writer
	err = data.ToJSON(finalKostList, rw)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)

		return
	}

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
	myKost, err := kostHandler.kost.GetKostByOwner(currentUser.ID)
	if err != nil {
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
	kostList, err := kostHandler.kost.GetKostListByOwner(currentUser.ID)
	if err != nil {
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

	geoLocation, err := kostHandler.kost.GetReverseGeocoderResult(userReq.Latitude, userReq.Longitude)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)

		return
	}

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
