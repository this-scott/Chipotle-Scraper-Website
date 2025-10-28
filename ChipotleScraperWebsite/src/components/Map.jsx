import { useEffect, useRef, useState } from "react";
import leaflet from "leaflet";

export default function Map() {
    //this creates the map as a DOM element
    const mapRef = useRef();
    const positionMarkerRefs = useRef({})

    const [data, setData] = useState([]);

    const decoder = new TextDecoder('utf-8');

    //This says the map exists even if it isn't rendered yet
    useEffect(() => {
        const fetchChipotles = (center) => {
            fetch(`http://localhost:8080/givechipotles?lat=${center.lat}&long=${center.lng}`)
                .then(response => response.json())
                .then(data => {
                    //YEAH NO NEED TO ITERATE WHEN YOU SEND IT SUCCESSFULLY

                    //creating a marker for each location in data. Data is going to need a pricemap soon as well :\
                    // storing everything in an array then passing it to state
                    // var chipotles = []
                    // data.forEach(item => {
                    //     //Just the command below would be nice if javascript was a good language
                    //     //leaflet.marker([userPosition.latitude, userPosition.longitude]).addTo(mapRef.current).bindPopup("Chipotle Information")
                    //     chipotles.push(item)
                    // })
                    setData(data)
                })
                .catch(error => {
                    console.error('Error fetching map data:', error);
                });
        }
        mapRef.current = leaflet.map('map').setView([37, -98.5795], 5);

        leaflet.tileLayer('https://tile.openstreetmap.org/{z}/{x}/{y}.png', {
            maxZoom: 19,
            attribution: '&copy; <a href="http://www.openstreetmap.org/copyright">OpenStreetMap</a>'
        }).addTo(mapRef.current);

        //move listener
        mapRef.current.on('move', () => {
            if (mapRef.current.getZoom() >= 10) {
                fetchChipotles(mapRef.current.getCenter());
            }
        })

        //zoom listener
        mapRef.current.on('zoomend', () => {
            if (mapRef.current.getZoom() >= 10) {
                fetchChipotles(mapRef.current.getCenter());
            }
          });

    },[]);

    //marker array attempt 1
    // useEffect(() => {
    //     if (positionMarkerRef.current) {
    //         mapRef.current.removeLayer(positionMarkerRef.current)
    //     }

    //     positionMarkerRef = leaflet.marker([])
    // }, [data])

    //OK WERE HANDLING MARKER ARRAY CHANGES HERE
    useEffect(() => {
        //convert bytes to json 
        // const jsonString = decoder.decode(byteArray);
        // const jsonObject = JSON
        //check the list of new markers to add the new ones
        //this is how react likes to iterate(fr this is so weird)
        data.forEach((pos) => {
            const {id, menu, latitude, longitude} = pos;

            if (!positionMarkerRefs.current[id]) {
                console.log({latitude, longitude})
                const marker = leaflet.marker([latitude, longitude])
                    .addTo(mapRef.current)
                    //DOM element. Either full component or rendered via a js function(seems easier for this case)
                    .bindPopup(<></>);

                positionMarkerRefs.current[id] = marker;
            }
        })

        //check the existing list to remove old ones
        Object.keys(positionMarkerRefs.current).forEach(id => {
            if (!data.find((p) => p.id == id)) {
                mapRef.current.removeLayer(positionMarkerRefs.current[id]);
                delete positionMarkerRefs.current[id];
            }
        })
    }, [data]);

    return  <div id="map" ref ={mapRef} ></div>;

}