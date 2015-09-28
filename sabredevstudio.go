package sabredevstudio

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"golang.org/x/net/context"
	"golang.org/x/oauth2/clientcredentials"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"strconv"
)

type DevStudioApiClient struct {
	*http.Client
	BaseUrl string
}

type Links struct {
	Links []struct {
		LinkURL string `json:"href"`
		Rel     string `json:"rel"`
	}
}
type Themes struct {
	Links
	Themes []struct {
		Links
		Theme string
	}
}
type Currency struct {
	AmountRaw     interface{} `json:"Amount"` // Not always float64 in the data
	Amount        float64     `json:",omitempty"`
	CurrencyCode  string
	DecimalPlaces int
	TaxCode       string `json:",omitempty"` // only applicable to taxes
}

// avoid recursion for the type cleansing dance below
type currency Currency

func (c *Currency) UnmarshalJSON(value []byte) error {
	var cleanedCurrency currency
	if err := json.Unmarshal(value, &cleanedCurrency); err != nil {
		panic(err)
	}
	if str, ok := cleanedCurrency.AmountRaw.(string); ok {
		cleanedCurrency.Amount, _ = strconv.ParseFloat(str, 64)
	} else if flt, ok := cleanedCurrency.AmountRaw.(float64); ok {
		cleanedCurrency.Amount = flt
	}
	*c = Currency{cleanedCurrency.AmountRaw, cleanedCurrency.Amount, cleanedCurrency.CurrencyCode, cleanedCurrency.DecimalPlaces, cleanedCurrency.TaxCode}
	return nil
}

type Equipment struct {
	AirEquipTypeRaw interface{} `json:"AirEquipType"`
	AirEquipType    string
}

// avoid recursion
type equipment Equipment

func (e *Equipment) UnmarshalJSON(value []byte) error {
	var cleanedEquipment equipment
	if err := json.Unmarshal(value, &cleanedEquipment); err != nil {
		panic(err)
	}
	if str, ok := cleanedEquipment.AirEquipTypeRaw.(string); ok {
		cleanedEquipment.AirEquipType = str
	} else if i, ok := cleanedEquipment.AirEquipTypeRaw.(int); ok {
		cleanedEquipment.AirEquipType = strconv.Itoa(i)
	}
	*e = Equipment{cleanedEquipment.AirEquipTypeRaw, cleanedEquipment.AirEquipType}
	return nil
}

type OnTimePerformance struct {
	LevelRaw interface{} `json:"Level"`
	Level    string
}

// avoid recursion
type ontimeperformance OnTimePerformance

func (o *OnTimePerformance) UnmarshalJSON(value []byte) error {
	var cleaned ontimeperformance
	if err := json.Unmarshal(value, &cleaned); err != nil {
		panic(err)
	}
	if str, ok := cleaned.LevelRaw.(string); ok {
		cleaned.Level = str
	} else if i, ok := cleaned.LevelRaw.(int); ok {
		cleaned.Level = strconv.Itoa(i)
	}
	*o = OnTimePerformance{cleaned.LevelRaw, cleaned.Level}
	return nil
}

type FlightShop struct {
	DepartureDateTime   string // Just the date
	ReturnDateTime      string // Just the date
	DestinationLocation string
	OriginLocation      string
	PricedItineraries   []struct {
		AirItinerary struct {
			DirectionInd             string
			OriginDestinationOptions struct {
				OriginDestinationOption []struct {
					ElapsedTime   int
					FlightSegment []struct {
						ArrivalAirport    struct{ LocationCode string }
						ArrivalDateTime   string // This one is the full timestamp
						ArrivalTimeZone   struct{ GMTOffset int }
						DepartureAirport  struct{ LocationCode string }
						DepartureDateTime string // This one is the full timestamp
						DepartureTimeZone struct{ GMTOffset int }
						ElapsedTime       int
						Equipment         Equipment
						FlightNumber      int
						MarketingAirline  struct{ Code string }
						MarriageGrp       string
						OnTimePerformance OnTimePerformance
						OperatingAirline  struct {
							FlightNumber int
							Code         string
						}
						ResBookDesignCode string
						StopQuantity      int
						TPA_Extensions    struct {
							eTicket struct{ Ind bool }
						}
					}
				}
			}
		}
		AirItineraryPricingInfo struct {
			AlternateCityOption bool
			FareInfos           struct {
				FareInfo []struct {
					FareReference  string
					TPA_Extensions struct {
						Cabin          struct{ Cabin string }
						SeatsRemaining struct {
							BelowMin bool
							Number   int
						}
					}
				}
			}
			ItinTotalFare struct {
				BaseFare         Currency
				EquivFare        Currency
				FareConstruction Currency
				Taxes            struct{ Tax []Currency }
				TotalFare        Currency
			}
			PTC_FareBreakdowns struct {
				PTC_FareBreakdown struct {
					FareBasisCodes struct {
						FareBasisCode []struct {
							ArrivalAirportCode   string
							AvailabilityBreak    bool
							BookingCode          string
							DepartureAirportCode string
							content              string
						}
					}
					PassengerFare struct {
						BaseFare         Currency
						EquivFare        Currency
						FareConstruction Currency
						TotalFare        Currency
					}
					PassengerTypeQuantity struct {
						Quantity int
						Code     string
					}
				}
				TPA_Extensions struct {
					DivideInParty struct{ Indicator bool }
				}
			}
		}
		AlternateAirport bool
		SequenceNumber   int
		TPA_Extensions   struct {
			ValidatingCarrier struct{ Code string }
		}
		TicketingInfo struct{ TicketType string }
	}
	Links
}

func clientID() (string, error) {
	clientID := os.Getenv("CLIENT_ID")
	clientID = base64.StdEncoding.EncodeToString([]byte(clientID))
	return clientID, nil
}
func clientSecret() (string, error) {
	clientSecret := os.Getenv("CLIENT_SECRET")
	clientSecret = base64.StdEncoding.EncodeToString([]byte(clientSecret))
	return clientSecret, nil
}
func baseUrl() (string, error) {
	baseUrl := os.Getenv("URL")
	return baseUrl, nil
}

func NewClient() *DevStudioApiClient {
	// Shout out to https://www.snip2code.com/Snippet/551369/Example-usage-of-https---godoc-org-golan
	baseUrl, _ := baseUrl()
	clientID, _ := clientID()
	clientSecret, _ := clientSecret()
	config := clientcredentials.Config{
		ClientID:     clientID,
		ClientSecret: clientSecret,
		TokenURL:     baseUrl + "/v1/auth/token",
	}
	// the client will update its token if it's expired
	client := config.Client(context.Background())
	return &DevStudioApiClient{Client: client, BaseUrl: baseUrl}
}
func (c *DevStudioApiClient) Request(requestUrl string) []byte {
	fmt.Printf("+%v\n", requestUrl)
	resp, err := c.Get(requestUrl)
	defer resp.Body.Close()
	if err != nil {
		panic(err)
	}
	content, _ := ioutil.ReadAll(resp.Body)
	return content
}
func (c *DevStudioApiClient) RequestWithParams(requestUrl string, params map[string]string) []byte {
	q := url.Values{}
	for key, value := range params {
		q.Add(key, value)
	}
	requestUrl = requestUrl + "?" + q.Encode()
	return c.Request(requestUrl)
}
func prettyPrintJson(content []byte) {
	var f interface{}
	_ = json.Unmarshal(content, &f)
	prettyJSON, _ := json.MarshalIndent(f, "", "  ")
	os.Stdout.Write(prettyJSON)
}

func (c *DevStudioApiClient) GetTravelThemes() (*Themes, error) {
	travelThemesUrl := c.BaseUrl + "/v1/lists/supported/shop/themes"
	content := c.Request(travelThemesUrl)
	//prettyPrintJson(content)
	var themes Themes
	if err := json.Unmarshal(content, &themes); err != nil {
		return nil, err
	}
	//fmt.Printf("+%v\n", themes)
	return &themes, nil
}
func (c *DevStudioApiClient) GetFlightSearch(params map[string]string) (*FlightShop, error) {
	flightSearchUrl := c.BaseUrl + "/v1/shop/flights"
	content := c.RequestWithParams(flightSearchUrl, params)
	//prettyPrintJson(content)
	var flightShop FlightShop
	if err := json.Unmarshal(content, &flightShop); err != nil {
		return nil, err
	}
	//fmt.Printf("+%v\n", flightShop)
	return &flightShop, nil
}
