package handlers

import (
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/fakhripraya/kost-service/config"
	"github.com/fakhripraya/kost-service/data"
	"github.com/fakhripraya/kost-service/database"
	"github.com/fakhripraya/kost-service/entities"
	"gorm.io/gorm"
)

// AddKost is a method to add the new given kost info to the database
func (kostHandler *KostHandler) AddKost(rw http.ResponseWriter, r *http.Request) {

	// get the kost via context
	kostReq := r.Context().Value(KeyKost{}).(*entities.Kost)

	// get the current user login
	var currentUser *database.MasterUser
	currentUser, err := kostHandler.kost.GetCurrentUser(rw, r, kostHandler.store)
	if err != nil {
		data.ToJSON(&GenericError{Message: err.Error()}, rw)

		return
	}

	// proceed to create the new kost with transaction scope
	err = config.DB.Transaction(func(tx *gorm.DB) error {

		// set variables
		var newKost database.DBKost
		var dbErr error

		newKost.OwnerID = currentUser.ID
		newKost.TypeID = kostReq.TypeID // kategori kos kosan atau kontrakan atau dll
		newKost.Status = 0              // TODO: create a documented status later // status 0 = baru
		newKost.KostCode, dbErr = kostHandler.kost.GenerateCode("K", kostReq.Country[0:1], kostReq.City[0:1])

		if dbErr != nil {
			return dbErr
		}

		newKost.KostName = kostReq.KostName
		newKost.KostDesc = kostReq.KostDesc
		newKost.Country = kostReq.Country
		newKost.City = kostReq.City
		newKost.Address = kostReq.Address
		newKost.Latitude = kostReq.Latitude
		newKost.Longitude = kostReq.Longitude
		newKost.ThumbnailURL = kostReq.Rooms[0].RoomPicts[0].URL
		newKost.UpRate = 0
		newKost.UpRateExpired = time.Now().Local()
		newKost.IsVerified = false
		newKost.IsActive = true
		newKost.Created = time.Now().Local()
		newKost.CreatedBy = currentUser.Username
		newKost.Modified = time.Now().Local()
		newKost.ModifiedBy = currentUser.Username

		if dbErr = tx.Create(&newKost).Error; dbErr != nil {
			return dbErr
		}

		// proceed to create the new kost periods with transaction scope
		dbErr = tx.Transaction(func(tx2 *gorm.DB) error {

			// create the variable specific to the nested transaction
			var dbErr2 error
			var kostPeriod = kostReq.KostPeriods

			// add the kost id to the slices
			for i := range kostPeriod {
				(&kostPeriod[i]).KostID = newKost.ID
				(&kostPeriod[i]).IsActive = true
				(&kostPeriod[i]).Created = time.Now().Local()
				(&kostPeriod[i]).CreatedBy = currentUser.Username
				(&kostPeriod[i]).Modified = time.Now().Local()
				(&kostPeriod[i]).ModifiedBy = currentUser.Username
			}

			// insert the new kost periods to database
			if dbErr2 = tx2.Create(&kostPeriod).Error; dbErr2 != nil {
				return dbErr2
			}

			// return nil will commit the whole nested transaction
			return nil
		})

		// if transaction error, return the error
		if dbErr != nil {
			return dbErr
		}

		// proceed to create the new kost picts with transaction scope
		dbErr = tx.Transaction(func(tx2 *gorm.DB) error {

			// create the variable specific to the nested transaction
			var dbErr2 error
			var kostPicts = kostReq.KostPicts

			// add the kost id to the slices
			for i := range kostPicts {
				(&kostPicts[i]).KostID = newKost.ID
				(&kostPicts[i]).IsActive = true
				(&kostPicts[i]).Created = time.Now().Local()
				(&kostPicts[i]).CreatedBy = currentUser.Username
				(&kostPicts[i]).Modified = time.Now().Local()
				(&kostPicts[i]).ModifiedBy = currentUser.Username
			}

			// insert the new kost picts to database
			if dbErr2 = tx2.Create(&kostPicts).Error; dbErr2 != nil {
				return dbErr2
			}

			// return nil will commit the whole nested transaction
			return nil
		})

		// if transaction error, return the error
		if dbErr != nil {
			return dbErr
		}

		// proceed to create the new kost benchmark with transaction scope
		dbErr = tx.Transaction(func(tx2 *gorm.DB) error {

			// create the variable specific to the nested transaction
			var dbErr2 error
			var kostBenchmark = kostReq.KostBenchmark

			// add the kost id to the slices
			for i := range kostBenchmark {
				(&kostBenchmark[i]).KostID = newKost.ID
				(&kostBenchmark[i]).IsActive = true
				(&kostBenchmark[i]).Created = time.Now().Local()
				(&kostBenchmark[i]).CreatedBy = currentUser.Username
				(&kostBenchmark[i]).Modified = time.Now().Local()
				(&kostBenchmark[i]).ModifiedBy = currentUser.Username
			}

			// insert the new kost picts to database
			if dbErr2 = tx2.Create(&kostBenchmark).Error; dbErr2 != nil {
				return dbErr2
			}

			// return nil will commit the whole nested transaction
			return nil
		})

		// if transaction error, return the error
		if dbErr != nil {
			return dbErr
		}

		// proceed to create the new kost accessibility with transaction scope
		dbErr = tx.Transaction(func(tx2 *gorm.DB) error {

			// create the variable specific to the nested transaction
			var dbErr2 error
			var kostAccess = kostReq.KostAccess

			// add the kost id to the slices
			for i := range kostAccess {
				(&kostAccess[i]).KostID = newKost.ID
				(&kostAccess[i]).IsActive = true
				(&kostAccess[i]).Created = time.Now().Local()
				(&kostAccess[i]).CreatedBy = currentUser.Username
				(&kostAccess[i]).Modified = time.Now().Local()
				(&kostAccess[i]).ModifiedBy = currentUser.Username
			}

			// insert the new kost picts to database
			if dbErr2 = tx2.Create(&kostAccess).Error; dbErr2 != nil {
				return dbErr2
			}

			// return nil will commit the whole nested transaction
			return nil
		})

		// if transaction error, return the error
		if dbErr != nil {
			return dbErr
		}

		// proceed to create the new around kost with transaction scope
		dbErr = tx.Transaction(func(tx2 *gorm.DB) error {

			// create the variable specific to the nested transaction
			var dbErr2 error
			var kostAround = kostReq.KostAround

			// add the kost id to the slices
			for i := range kostAround {
				(&kostAround[i]).KostID = newKost.ID
				(&kostAround[i]).IsActive = true
				(&kostAround[i]).Created = time.Now().Local()
				(&kostAround[i]).CreatedBy = currentUser.Username
				(&kostAround[i]).Modified = time.Now().Local()
				(&kostAround[i]).ModifiedBy = currentUser.Username
			}

			// insert the new kost picts to database
			if dbErr2 = tx2.Create(&kostAround).Error; dbErr2 != nil {
				return dbErr2
			}

			// return nil will commit the whole nested transaction
			return nil
		})

		// if transaction error, return the error
		if dbErr != nil {
			return dbErr
		}

		// loop the room slices
		for _, room := range kostReq.Rooms {

			// add the kostReq room slices into the database
			dbErr = kostHandler.kost.AddRoom(currentUser, newKost.ID, &room)

			// if transaction error, return the error
			if dbErr != nil {
				return dbErr
			}

		}

		// add the kost facilities to the database
		dbErr = kostHandler.kost.AddFacilities(currentUser, newKost.ID, kostReq.Facilities)

		// if transaction error, return the error
		if dbErr != nil {
			return dbErr
		}

		// return nil will commit the whole transaction
		return nil

	})

	// if transaction error
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)

		return
	}

	rw.WriteHeader(http.StatusOK)
	data.ToJSON(&GenericError{Message: "Sukses menambah kost baru"}, rw)
	return
}

func (kostHandler *KostHandler) AddKostAds(rw http.ResponseWriter, r *http.Request) {

	// get the kost via context
	kostAdsReq := r.Context().Value(KeyKostAds{}).(*entities.KostAds)

	// proceed to create the new kost ads with transaction scope
	err := config.DB.Transaction(func(tx *gorm.DB) error {

		// set variables
		var newKostAds database.DBKostAds
		var dbErr error

		newKostAds.Status = 0
		newKostAds.AdsCode, dbErr = kostHandler.kost.GenerateCode("K", "I", "ADS")

		if dbErr != nil {
			return dbErr
		}

		newKostAds.AdsType = kostAdsReq.AdsType
		newKostAds.AdsKostType = kostAdsReq.AdsKostType
		newKostAds.AdsOwner = kostAdsReq.AdsOwner
		newKostAds.AdsOwnerIG = kostAdsReq.AdsOwnerIG
		newKostAds.AdsPhoneNumber = kostAdsReq.AdsPhoneNumber
		newKostAds.AdsPICWhatsapp = kostAdsReq.AdsPICWhatsapp
		newKostAds.AdsPropertyAddress = kostAdsReq.AdsPropertyAddress
		newKostAds.AdsPropertyCity = kostAdsReq.AdsPropertyCity
		newKostAds.AdsPropertyPrice = kostAdsReq.AdsPropertyPrice
		newKostAds.AdsDesc = kostAdsReq.AdsDesc
		newKostAds.AdsGender = kostAdsReq.AdsGender
		newKostAds.AdsPetAllowed = kostAdsReq.AdsPetAllowed
		newKostAds.AdsPostScheduleRequest = kostAdsReq.AdsPostScheduleRequest
		newKostAds.AdsHashtag = kostAdsReq.AdsHashtag
		newKostAds.AdsLinkSwipeUp = kostAdsReq.AdsLinkSwipeUp
		newKostAds.AdsIgBioLink = kostAdsReq.AdsIgBioLink
		newKostAds.IsActive = true
		newKostAds.Created = time.Now().Local()
		newKostAds.CreatedBy = "System"
		newKostAds.Modified = time.Now().Local()
		newKostAds.ModifiedBy = "System"

		if dbErr = tx.Create(&newKostAds).Error; dbErr != nil {
			return dbErr
		}

		var baseFileDirPath = os.Getenv("APP_MK_DIR_ADS_FILE_PATH") + "/File-Directory-" + strconv.FormatUint(uint64(newKostAds.ID), 10)

		// the WriteFile method returns an error if unsuccessful
		dbErr = os.MkdirAll(baseFileDirPath, 0777)
		if dbErr != nil {
			return dbErr
		}

		// proceed to create the new kost ads files with transaction scope
		dbErr = tx.Transaction(func(tx2 *gorm.DB) error {

			// create the variable specific to the nested transaction
			var dbErr2 error
			var kostAdsFiles = kostAdsReq.AdsFiles

			// add the kost id to the slices
			for i := range kostAdsFiles {

				// specify the directory path of the desired file
				var specificDirPath = baseFileDirPath + "/FOLDER_" + (&kostAdsFiles[i]).AdsFileType

				(&kostAdsFiles[i]).AdsID = newKostAds.ID
				(&kostAdsFiles[i]).AdsDirPath = specificDirPath
				(&kostAdsFiles[i]).IsActive = true
				(&kostAdsFiles[i]).Created = time.Now().Local()
				(&kostAdsFiles[i]).CreatedBy = "System"
				(&kostAdsFiles[i]).Modified = time.Now().Local()
				(&kostAdsFiles[i]).ModifiedBy = "System"

				_, err := os.Stat(specificDirPath)
				if os.IsNotExist(err) {
					// the WriteFile method returns an error if unsuccessful
					dbErr = os.Mkdir(specificDirPath, 0777)
					if dbErr != nil {
						return dbErr
					}
				}

				dbErr = ioutil.WriteFile(specificDirPath+"/"+strconv.Itoa(i)+".txt", []byte((&kostAdsFiles[i]).BASE64STRING), 0777)
				// handle this error
				if dbErr != nil {
					// print it out
					return dbErr
				}

				(&kostAdsFiles[i]).BASE64STRING = ""
			}

			// insert the new kost periods to database
			if dbErr2 = tx2.Create(&kostAdsFiles).Error; dbErr2 != nil {
				return dbErr2
			}

			//return nil will commit the whole nested transaction
			return nil
		})

		// if transaction error, return the error
		if dbErr != nil {
			return dbErr
		}

		// return nil will commit the whole transaction
		return nil

	})

	// if transaction error
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		data.ToJSON(&GenericError{Message: err.Error()}, rw)

		return
	}

	rw.WriteHeader(http.StatusOK)
	data.ToJSON(&GenericError{Message: "Sukses submit form iklan, sekarang kamu hanya tinggal tunggu kami proses deh, kalau menurut kamu kelamaan dipostnya jangan lupa untuk tegur kita ya :)"}, rw)
	return
}
