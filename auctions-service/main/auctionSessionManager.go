package main

import (
	"time"
)

const (
	loadAheadDuration  time.Duration = time.Duration(2) * time.Hour // how much in advance AuctionSessionManager should bring Auctions into memory before their start
	loadBehindDuration time.Duration = time.Duration(2) * time.Hour // how much in the past AuctionSessionManager should load Auctions into memory since their end
	// note, with a system shutdown and restart, AuctionSessionManager may have to finalize Auctions that ended e.g. 30 minutes ago
	FinalizeDelay time.Duration = time.Duration(30) * time.Minute // how long after auction end before SessionManager should finalize the auction
	// note, with a system shutdown and restart, the context may want to process bids sitting in queue that were submitted during live auction
	// session. so, give this context 30 minutes to get up to speed with its bid processing, which may involve processing bids as if from an
	// earlier point in time.
)

type AuctionSessionManager struct {
	auctionsservice  *AuctionService
	turnedOn         bool
	lastAlertTime    time.Time
	lastFinalizeTime time.Time
	lastLoadTime     time.Time
	alertCycle       time.Duration
	finalizeCycle    time.Duration
	loadCycle        time.Duration
}

func NewAuctionSessionManager(
	auctionsservice *AuctionService,
	alertCycle time.Duration,
	finalizeCycle time.Duration,
	loadAuctionCycle time.Duration,
) *AuctionSessionManager {

	lastAlertTime := time.Now()
	lastFinalizeTime := time.Now()
	lastLoadTime := time.Now()

	return &AuctionSessionManager{
		auctionsservice,
		false,
		lastAlertTime,
		lastFinalizeTime,
		lastLoadTime,
		alertCycle,
		finalizeCycle,
		loadAuctionCycle,
	}
}

func (auctionSessionManager *AuctionSessionManager) TurnOn() {
	if !auctionSessionManager.turnedOn {
		auctionSessionManager.turnedOn = true
		auctionSessionManager.lastAlertTime = time.Now()
		auctionSessionManager.lastFinalizeTime = time.Now()
		auctionSessionManager.lastLoadTime = time.Now()

		// load into memory almost all past auctions because we don't know how long the server has been down.
		// might need to finalize very old auctions that are over but have not been concluded and archived.
		// load auctions whose start->end period overlap with the time period from jan 1, 1950 to ~2 hrs
		// ahead of present moment. this will load in the auctions that'll start <2 hrs from now.
		since := time.Date(1950, 1, 1, 0, 00, 00, 0, time.UTC) // load from jan 1, 1950
		upTo := time.Now().Add(loadAheadDuration)              // up to ~ 2hrs from now

		auctionSessionManager.auctionsservice.LoadAuctionsIntoMemory(since, upTo)
		auctionSessionManager.auctionsservice.SendOutLifeCycleAlerts()
		auctionSessionManager.auctionsservice.FinalizeAnyPastAuctions(FinalizeDelay)

		go auctionSessionManager.intermittentlyLoadAuctions()
		go auctionSessionManager.intermittentlySendLifeCycleAlerts()
		go auctionSessionManager.intermittentlyFinalizeAuctions()
	}
}

func (auctionSessionManager *AuctionSessionManager) TurnOff() {
	if auctionSessionManager.turnedOn {
		auctionSessionManager.turnedOn = false // this will terminate 3 asynch goroutines
	}
}

func (auctionSessionManager *AuctionSessionManager) intermittentlyLoadAuctions() {
	for auctionSessionManager.turnedOn {
		if time.Since(auctionSessionManager.lastLoadTime) >= auctionSessionManager.loadCycle {
			since := auctionSessionManager.lastLoadTime.Add(loadAheadDuration)
			upTo := time.Now().Add(loadAheadDuration)
			auctionSessionManager.auctionsservice.LoadAuctionsIntoMemory(since, upTo) // acquires lock
			auctionSessionManager.lastLoadTime = time.Now()
		}
	}
}

func (auctionSessionManager *AuctionSessionManager) intermittentlySendLifeCycleAlerts() {
	for auctionSessionManager.turnedOn {
		if time.Since(auctionSessionManager.lastAlertTime) >= auctionSessionManager.alertCycle {
			// since := auctionSessionManager.lastAlertTime.Add(loadAheadDuration)
			// upTo := time.Now().Add(loadAheadDuration)
			auctionSessionManager.auctionsservice.SendOutLifeCycleAlerts() // acquires lock
			auctionSessionManager.lastAlertTime = time.Now()
		}
	}
}

func (auctionSessionManager *AuctionSessionManager) intermittentlyFinalizeAuctions() {
	for auctionSessionManager.turnedOn {
		if time.Since(auctionSessionManager.lastFinalizeTime) >= auctionSessionManager.finalizeCycle {
			// since := auctionSessionManager.lastLoadTime.Add(loadAheadDuration)
			// upTo := time.Now().Add(loadAheadDuration)
			auctionSessionManager.auctionsservice.FinalizeAnyPastAuctions(FinalizeDelay) // acquires lock
			auctionSessionManager.lastFinalizeTime = time.Now()
		}
	}
}
