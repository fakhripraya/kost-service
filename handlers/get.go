package handlers

import (
	"io/ioutil"
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
	kostList, _, err := kostHandler.kost.GetKostListByOwner(kostOwner.ID, -1)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)

		return
	}

	finalKostAround := struct {
		ID             uint            `json:"id"`
		DisplayName    string          `json:"display_name"`
		ProfilePicture string          `json:"profile_picture"`
		City           string          `json:"city"`
		KostList       []entities.Kost `json:"kost_list"`
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

	// get the room id via mux
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

// GetKostRoomInfoAll is a method to fetch the given kost room detail list with kost id as the parameter
func (kostHandler *KostHandler) GetKostRoomInfoAll(rw http.ResponseWriter, r *http.Request) {

	// get the kost via mux
	vars := mux.Vars(r)
	page, err := strconv.Atoi(vars["page"])
	kostID, err := strconv.ParseUint(vars["id"], 10, 32)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		data.ToJSON(&GenericError{Message: "Unable to convert id"}, rw)

		return
	}
	kostRoomDetails, err := kostHandler.kost.GetKostRoomDetailsByKost(uint(kostID), page)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)

		return
	}

	var kostRoomDetailSum []database.DBKostRoomDetail
	if err := config.DB.Where("kost_id = ?", kostID).Find(&kostRoomDetailSum).Error; err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)

		return
	}

	var kostRoomDetailsFinal []entities.KostRoomDetail
	for index, roomDetail := range kostRoomDetails {

		if index >= ((page * 10) - 10) {

			kostRoom, err := kostHandler.kost.GetKostRoom(roomDetail.RoomID)
			if err != nil {
				rw.WriteHeader(http.StatusBadRequest)
				data.ToJSON(&GenericError{Message: err.Error()}, rw)

				return
			}

			currency, err := kostHandler.kost.GetUOMDesc(kostRoom.RoomPriceUOM)
			if err != nil {
				rw.WriteHeader(http.StatusBadRequest)
				data.ToJSON(&GenericError{Message: err.Error()}, rw)

				return
			}

			kostRoomDetailBook, err := kostHandler.kost.GetKostRoomBooked(roomDetail.ID)
			if err != nil {
				rw.WriteHeader(http.StatusBadRequest)
				data.ToJSON(&GenericError{Message: err.Error()}, rw)

				return
			}

			if kostRoomDetailBook != nil {
				period, err := kostHandler.kost.GetMasterPeriod(kostRoomDetailBook.PeriodID)
				if err != nil {
					rw.WriteHeader(http.StatusBadRequest)
					data.ToJSON(&GenericError{Message: err.Error()}, rw)

					return
				}

				booker, err := kostHandler.kost.GetMasterUser(kostRoomDetailBook.BookerID)
				if err != nil {
					rw.WriteHeader(http.StatusBadRequest)
					data.ToJSON(&GenericError{Message: err.Error()}, rw)

					return
				}

				// TODO: status ntr diomongin lagi mau gimana
				kostRoomDetailsFinal = append(kostRoomDetailsFinal, entities.KostRoomDetail{
					ID:         roomDetail.ID,
					KostID:     roomDetail.KostID,
					RoomID:     roomDetail.RoomID,
					RoomDesc:   kostRoom.RoomDesc,
					RoomNumber: roomDetail.RoomNumber,
					FloorLevel: roomDetail.FloorLevel,
					Price:      kostRoom.RoomPrice,
					Currency:   currency,
					Status:     kostRoomDetailBook.Status,
					Booker: &database.MasterUser{
						ID:             booker.ID,
						DisplayName:    booker.DisplayName,
						ProfilePicture: booker.ProfilePicture,
					},
					PrevPayment: kostRoomDetailBook.BookDate,
					NextPayment: kostRoomDetailBook.BookDate.AddDate(0, 0, int(period.PeriodValue)),
					IsActive:    roomDetail.IsActive,
				})
			} else {
				kostRoomDetailsFinal = append(kostRoomDetailsFinal, entities.KostRoomDetail{
					ID:         roomDetail.ID,
					KostID:     roomDetail.KostID,
					RoomID:     roomDetail.RoomID,
					RoomDesc:   kostRoom.RoomDesc,
					RoomNumber: roomDetail.RoomNumber,
					FloorLevel: roomDetail.FloorLevel,
					Price:      kostRoom.RoomPrice,
					Currency:   currency,
					IsActive:   roomDetail.IsActive,
				})
			}

		}

	}

	finalData := struct {
		KostRoomDetailSum int                       `json:"kost_room_detail_sum"`
		KostRoomDetail    []entities.KostRoomDetail `json:"kost_room_detail"`
	}{
		KostRoomDetailSum: len(kostRoomDetailSum),
		KostRoomDetail:    kostRoomDetailsFinal,
	}

	// parse the given instance to the response writer
	err = data.ToJSON(finalData, rw)
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
	// 6 = My kost list
	var count int64
	var kostList []entities.Kost
	if category == 0 {
		kostList, count, err = kostHandler.kost.GetKostList(page)
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

		kostList, count, err = kostHandler.kost.GetNearbyKostList(userReq.Latitude, userReq.Longitude, page)
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			data.ToJSON(&GenericError{Message: err.Error()}, rw)

			return
		}
	} else if category == 6 {
		// get the current user login
		var currentUser *database.MasterUser
		currentUser, err := kostHandler.kost.GetCurrentUser(rw, r, kostHandler.store)
		if err != nil {
			data.ToJSON(&GenericError{Message: err.Error()}, rw)

			return
		}

		// look for the current kost list in the db
		kostList, count, err = kostHandler.kost.GetKostListByOwner(currentUser.ID, page)
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

	type FinalResult struct {
		KostList  []FinalKostList `json:"kost_list"`
		KostCount int64           `json:"kost_count"`
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
	err = data.ToJSON(FinalResult{
		KostList:  finalKostList,
		KostCount: count,
	}, rw)
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
		OwnerID      uint    `json:"owner_id"`
	}

	// for this request, the page will always 2
	listNearbyKosts, _, err := kostHandler.kost.GetNearbyKostList(userReq.Latitude, userReq.Longitude, 2)
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
			OwnerID:      nearby.OwnerID,
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

// GetKostInstagramAdsList is a method to fetch the given kost Instagram ads list
func (kostHandler *KostHandler) GetKostInstagramAdsList(rw http.ResponseWriter, r *http.Request) {

	// get the kost via context
	//kostReq := r.Context().Value(KeyKostAds{}).(*entities.KostAds)

	var kostAds []entities.KostAds

	if err := config.DB.
		Model(&database.DBKostAds{}).
		Select("db_kost_ads.id"+
			",db_kost_ads.status "+
			",db_kost_ads.ads_code"+
			",db_kost_ads.ads_type"+
			",db_kost_ads.ads_kost_type"+
			",db_kost_ads.ads_owner"+
			",db_kost_ads.ads_owner_ig"+
			",db_kost_ads.ads_phone_number"+
			",db_kost_ads.ads_pic_whatsapp"+
			",db_kost_ads.ads_property_address"+
			",db_kost_ads.ads_property_city"+
			",db_kost_ads.ads_property_price"+
			",db_kost_ads.ads_desc"+
			",db_kost_ads.ads_gender"+
			",db_kost_ads.ads_pet_allowed"+
			",db_kost_ads.ads_post_schedule_request"+
			",db_kost_ads.ads_hashtag"+
			",db_kost_ads.is_active").
		Where("db_kost_ads.ads_type != ? OR db_kost_ads.ads_type != ?", "Iklan Tiktok Free", "Iklan Tiktok Premium (Rekomendasi)").Scan(&kostAds).Error; err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)

		return
	}

	// parse the given instance to the response writer
	err := data.ToJSON(kostAds, rw)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)

		return
	}

	return
}

// GetKostInstagramAdsFileList is a method to fetch the given kost Instagram ads file list
func (kostHandler *KostHandler) GetKostInstagramAdsFileList(rw http.ResponseWriter, r *http.Request) {

	// get the kost via context
	//kostReq := r.Context().Value(KeyKostAds{}).(*entities.KostAds)

	// get the id via mux
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)

		return
	}

	var kostAdsFiles []database.DBKostAdsFiles

	if err = config.DB.
		Model(&database.DBKostAdsFiles{}).
		Select("db_kost_ads_files.id"+
			",db_kost_ads_files.ads_id "+
			",db_kost_ads_files.ads_file_type"+
			",db_kost_ads_files.ads_dir_path"+
			",db_kost_ads_files.is_active").
		Where("db_kost_ads_files.ads_id = ?", id).Scan(&kostAdsFiles).Error; err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)

		return
	}

	// add the kost id to the slices
	for i := range kostAdsFiles {

		body, err := ioutil.ReadFile((&kostAdsFiles[i]).AdsDirPath)
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			data.ToJSON(&GenericError{Message: err.Error()}, rw)

			return
		}

		(&kostAdsFiles[i]).BASE64STRING = string(body)
	}

	// parse the given instance to the response writer
	err = data.ToJSON(kostAdsFiles, rw)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)

		return
	}

	return
}
