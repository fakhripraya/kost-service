package handlers

import (
	"net/http"
	"sort"
	"strconv"

	"github.com/fakhripraya/kost-service/config"
	"github.com/fakhripraya/kost-service/data"
	"github.com/fakhripraya/kost-service/database"
	"github.com/fakhripraya/kost-service/entities"
	"github.com/gorilla/mux"
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

// GetKostPeriod is a method to fetch the given kost period list
func (kostHandler *KostHandler) GetKostPeriod(rw http.ResponseWriter, r *http.Request) {

	// get the kost via context
	kostReq := r.Context().Value(KeyKost{}).(*entities.Kost)

	var kostPeriod []entities.KostPeriod
	if err := config.DB.
		Model(&database.DBKostPeriod{}).
		Select("db_kost_periods.id,db_kost_periods.kost_id,db_kost_periods.period_id,master_periods.period_desc, master_periods.is_active").
		Joins("inner join master_periods on master_periods.id = db_kost_periods.period_id").
		Where("db_kost_periods.kost_id = ?", kostReq.ID).Scan(&kostPeriod).Error; err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)

		return
	}

	// parse the given instance to the response writer
	err := data.ToJSON(kostPeriod, rw)
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

	// get the additional field
	vars := mux.Vars(r)

	// get the facilities from the method
	kostFacilities, kostRoomFacilities, err := kostHandler.kost.GetKostFacilities(kostReq.ID, vars["roomId"])
	if err != nil {
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
		Where("db_kost_reviews.kost_id = ?", kostReq.ID).Scan(&kostReview).Error; err != nil {
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

	// get the current logged in user additional info via context
	userReq := r.Context().Value(KeyUser{}).(*entities.User)

	// get the page via mux
	vars := mux.Vars(r)
	page, err := strconv.Atoi(vars["page"])
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)

		return
	}

	category, err := strconv.Atoi(vars["category"])
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)

		return
	}

	// 0 = all kost // Initial val
	// 1 = Near You
	//TODO: 2 = Most Popular
	//TODO: 3 = Most Facilited
	//TODO: 4 = Most Expensive
	//TODO: 5 = Most Cheap
	var kostList []entities.Kost
	if category == 0 {
		kostList, err = kostHandler.kost.GetKostList(page)
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			data.ToJSON(&GenericError{Message: err.Error()}, rw)

			return
		}
	} else if category == 1 {
		if userReq.Latitude == "" || userReq.Longitude == "" {
			// if latitude or longitude is an empty string
			// parse the given instance to the response writer
			err = data.ToJSON(kostList, rw)
			if err != nil {
				rw.WriteHeader(http.StatusBadRequest)
				data.ToJSON(&GenericError{Message: err.Error()}, rw)

				return
			}
			return
		}

		kostList, err = kostHandler.kost.GetNearbyKostList(userReq.Latitude, userReq.Longitude, page)
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			data.ToJSON(&GenericError{Message: err.Error()}, rw)

			return
		}
	}

	// Sort by distance, keeping original order or equal elements.
	sort.SliceStable(kostList, func(i, j int) bool {
		return kostList[i].ID < kostList[j].ID
	})

	type FinalKostList struct {
		Kost       entities.Kost             `json:"kost"`
		Facilities []entities.KostFacilities `json:"facilities"`
		Price      float64                   `json:"price"`
		Currency   string                    `json:"currency"`
	}

	var finalKostList []FinalKostList
	for index, kost := range kostList {

		if index >= ((page * 10) - 10) {

			// get the facilities from the method
			kostFacilities, _, err := kostHandler.kost.GetKostFacilities(kost.ID, "")
			if err != nil {
				rw.WriteHeader(http.StatusBadRequest)
				data.ToJSON(&GenericError{Message: err.Error()}, rw)

				return
			}

			lowestPrice, err := kostHandler.kost.GetLowestPrice(kost.ID)
			if err != nil {
				rw.WriteHeader(http.StatusBadRequest)
				data.ToJSON(&GenericError{Message: err.Error()}, rw)

				return
			}

			finalKostList = append(finalKostList, FinalKostList{
				Kost:       kost,
				Facilities: kostFacilities,
				Price:      lowestPrice.RoomPrice,
				Currency:   lowestPrice.RoomPriceUomDesc,
			})
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

	type NearbyKostView struct {
		ID           uint    `json:"id"`
		KostName     string  `json:"kost_name"`
		City         string  `json:"city"`
		ThumbnailURL string  `json:"thumbnail_url"`
		Price        float64 `json:"price"`
		Currency     string  `json:"currency"`
	}

	// for this request, the page will always 2
	listNearbyKosts, err := kostHandler.kost.GetNearbyKostList(userReq.Latitude, userReq.Longitude, 2)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)

		return
	}

	// variable to hold temporary data of nearby kost list
	var finalNearbyKostList []NearbyKostView

	for _, nearby := range listNearbyKosts {

		lowestPrice, err := kostHandler.kost.GetLowestPrice(nearby.ID)
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			data.ToJSON(&GenericError{Message: err.Error()}, rw)

			return
		}

		finalNearbyKostList = append(finalNearbyKostList, NearbyKostView{
			ID:           nearby.ID,
			KostName:     nearby.KostName,
			City:         nearby.City,
			ThumbnailURL: nearby.ThumbnailURL,
			Price:        lowestPrice.RoomPrice,
			Currency:     lowestPrice.RoomPriceUomDesc,
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
