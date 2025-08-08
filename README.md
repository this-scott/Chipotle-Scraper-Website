# Chipotle Price Index

## As of right now this is a map displaying nearby chipotle locations

This is meant to be a continuation of my [Chipotle scraper project](https://github.com/this-scott/Chipotle-Scraper). The plan is to turn this into an interactive map displaying each stores menu on each marker and visuals to compare states and counties but I'm not sure long I'm going to stay on this project.

The dataset I created is in the backend directory.

## [Demo Gif of current state](https://scott-cv-presentation-slide-resources.s3.us-east-1.amazonaws.com/ChipotleMap.gif)
![ChipMap](https://scott-cv-presentation-slide-resources.s3.us-east-1.amazonaws.com/ChipotleMap.gif)

## How it works
* JSON containing Chipotle Locations and menu information is cached in memory of a golang endpoint.
* Vite + React stack serves a Leaflet map to the client in a component which holds the nearby positions as a state variable
* The Map is held as a reference
* The list of markers held as a reference
* Zoom and movement listeners on the map query the golang endpoint for Chipotles within 1DD(should probably use miles) of map center and writes the output to the positions state
* A change in the positions state triggers an effect which checks the new state to place unplaced markers, then checks the old list to find and remove markers which are not included in the new state

## TODO
* Consider websocket after http works
* Add menus
* Overlays comparing prices between states and counties in each state

Not affiliated with Chipotle. Go eat there, their food is really good. 
