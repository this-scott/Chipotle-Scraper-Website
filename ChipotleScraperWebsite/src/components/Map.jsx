import { useEffect, useRef, useState } from "react";
import leaflet from "leaflet";
import '../assets/markerformat.css';

export default function Map() {
    //this creates the map as a DOM element
    const mapRef = useRef();
    const positionMarkerRefs = useRef({});
    const chipGroupRef = useRef();
    const [data, setData] = useState([]);

    
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

        chipGroupRef.current = leaflet.layerGroup().addTo(mapRef.current);

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
            } else {
                chipGroupRef.current.clearLayers();
                positionMarkerRefs.current = {};                
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
                //console.log({latitude, longitude})
                //render menu right here
                // Create the grid container
                const newdiv = document.createElement('div');
                newdiv.classList.add('grid-container');
                
                let added = Set(['CHICKEN', 'STEAK', 'BEEF BARBACOA', 'CARNITAS', 'SOFRITAS', 'VEGGIE', 'CHIPS', ]);
                //add popular elements first(meats and mexican coke)
                const items = ['CHICKEN', 'STEAK', 'BEEF BARBACOA', 'CARNITAS', 'SOFRITAS', 'VEGGIE']

                // Create and append grid items
                Object.entries(menu).forEach(([key, value]) => {
                    console.log(key,value);
                    const gridItem = document.createElement('div');
                    gridItem.classList.add('grid-item');


                    const img = document.createElement('img');
                    img.src = item.img;
                    img.alt = item.title;

                    // Text (title + subtext)
                    const textContainer = document.createElement('div');
                    textContainer.classList.add('text-container');

                    const title = document.createElement('h3');
                    title.textContent = item.title;

                    const subtext = document.createElement('p');
                    subtext.textContent = item.subtext;

                    textContainer.appendChild(title);
                    textContainer.appendChild(subtext);

                    // Assemble grid item
                    gridItem.appendChild(img);
                    gridItem.appendChild(textContainer);

                    // Add to grid container
                    newdiv.appendChild(gridItem);
                });

                //for each element
                //div.innerHTML = 'Custom content';

                const marker = leaflet.marker([latitude, longitude])
                    .addTo(chipGroupRef.current)
                    //DOM element. Either full component or rendered via a js function(seems easier for this case)
                    .bindPopup(newdiv);

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