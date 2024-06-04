const map = L.map('map').setView([46.603354, 1.888334], 6);
L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
  attribution: '© OpenStreetMap contributors'
}).addTo(map);

const randomLocationBtn = document.getElementById('random-location-btn');
let hotelMarker = null;

const francePolygon = [
  [51.09931136914983, 2.5596690000205626],
  [50.91082422447218, 1.656558855932957],
  [50.194576243555076, 1.582314183391361],
  [49.55547153393911, 0.07110844621644219],
  [49.47603652483795, 0.4404656320381264],
  [49.384723873587575, -1.2253465142589448],
  [49.704483906124324, -1.3086814442448826],
  [49.71776429358721, -1.930146394721362],
  [48.604339826903725, -1.4562448874389986],
  [48.83938703333382, -3.1785853695628816],
  [48.50444268909851, -4.790532992670904],
  [48.32088049152992, -4.770845187938761],
  [48.295917504147326, -4.274283981637126],
  [48.20093006174537, -4.617092426441701],
  [48.124078035172175, -4.25806715040045],
  [48.00129647660796, -4.702536608599871],
  [47.473220774498145, -2.5179049071804513],
  [46.8190922670363, -2.1306411650147083],
  [46.11702141681366, -1.107300557093481],
  [44.50908697762128, -1.2030391867158698],
  [43.36188688093836, -1.6416633322764937],
  [42.71835564166801, 0.16965131109824938],
  [42.907940803823465, 0.8057849544499049],
  [42.51959085499763, 3.0392971135362643],
  [43.156937175932654, 3.046937304377991],
  [43.488211734862034, 4.063974597977307],
  [43.110780341404165, 6.402234465935635],
  [43.79061692414368, 7.373227174649713],
  [44.230220510017745, 7.592033048259225],
  [44.48221554674805, 6.902307944136567],
  [44.963991921594385, 6.934466901767678],
  [45.15617727329149, 6.567467901755265],
  [45.43446239255712, 7.110688774564551],
  [45.88730522006858, 6.9523095208091945],
  [46.34723006637046, 6.600035724964471],
  [46.00289483133537, 5.9998543200710515],
  [46.845349095173844, 6.0524633122788885],
  [47.217904503736435, 6.713370028850733],
  [47.64695837350456, 6.846091843460073],
  [47.68775842730403, 7.344373559490123],
  [48.832215122102724, 7.7578687745501895],
  [49.82068619792767, 4.485117256085061],
  [50.600194536315996, 2.97398358147467],
  [51.09931136914983, 2.5596690000205626]
];

function isPointInPolygon(point, vs) {
  const x = point[1], y = point[0];

  let inside = false;
  for (let i = 0, j = vs.length - 1; i < vs.length; j = i++) {
    const xi = vs[i][1], yi = vs[i][0];
    const xj = vs[j][1], yj = vs[j][0];

    const intersect = ((yi > y) != (yj > y)) &&
                      (x < (xj - xi) * (y - yi) / (yj - yi) + xi);
    if (intersect) inside = !inside;
  }

  return inside;
}

randomLocationBtn.addEventListener('click', () => {
  let randomLatLng;
  do {
    const lat = 42.51959085499763 + Math.random() * (51.09931136914983 - 42.51959085499763);
    const lng = -5.142222 + Math.random() * (8.231167 - -5.142222);
    randomLatLng = [lat, lng];
  } while (!isPointInPolygon(randomLatLng, francePolygon));

  if (hotelMarker) {
    map.removeLayer(hotelMarker);
  }

  hotelMarker = L.marker(randomLatLng).addTo(map)
    .bindPopup(`Vous avez découvert : [${randomLatLng[0].toFixed(5)}, ${randomLatLng[1].toFixed(5)}]`)
    .openPopup();

  map.setView(randomLatLng, 10);

  fetch('/update-coordinates', {
    method: 'POST',
    headers: {
      'Content-Type': 'application/json',
    },
    body: JSON.stringify({ latitude: randomLatLng[0], longitude: randomLatLng[1] }),
  })
  .then(response => response.json())
  .then(data => {
    console.log('Success:', data);
  })
  .catch((error) => {
    console.error('Error:', error);
  });

  fetch('/hotel')
    .then(response => response.json())
    .then(data => {
      const hotels = data.data.slice(0, 5);
      updateHotelInfo(hotels); 
    })
    .catch(error => {
      console.error('Error fetching hotel data:', error);
      updateHotelInfo([]); 
    });

    fetch('/activitie')
    .then(response => response.json())
    .then(data => {
      const activitie = data.data.slice(0, 5);
      updateActiviteInfo(activitie); 
    })
    .catch(error => {
      console.error('Error fetching activite data:', error);
      updateHotelInfo([]); 
    });
});

function updateActivitieInfo(activitie) {
  const activitieInfoContainer = document.getElementById("activitie-info");
  activitieInfoContainer.innerHTML = ""; 

  if (activitie.length === 0) {
    const noActivitieMessage = document.createElement('p');
    noActivitieMessage.textContent = "Aucune Activité trouvé.";
    activitieInfoContainer.appendChild(noActivitieMessage);
    return;
  }
  activitie.forEach(activitie => {
    const activitieContainer = document.createElement('div');
    activitieContainer.className = 'hotel';

    const activitieName = document.createElement('h2');
    activitieName.textContent = hotel.name;

    const activitieAddress = document.createElement('p');
    activitieAddress.textContent = `Country: ${activitie.address.countryCode}`;

    const activitieDistance = document.createElement('p');
    activitieDistance.textContent = `Distance: ${activitie.distance.value} ${activitie.distance.unit}`;

    activitieContainer.appendChild(activitieName);
    activitieContainer.appendChild(activitieAddress);
    activitieContainer.appendChild(activitieDistance);
    activitieInfoContainer.appendChild(activitieContainer);
  });
}


function updateHotelInfo(hotels) {
  const hotelInfoContainer = document.getElementById("hotel-info");
  hotelInfoContainer.innerHTML = ""; 

  if (hotels.length === 0) {
    const noHotelsMessage = document.createElement('p');
    noHotelsMessage.textContent = "Aucun hôtel trouvé.";
    hotelInfoContainer.appendChild(noHotelsMessage);
    return;
  }

  hotels.forEach(hotel => {
    const hotelContainer = document.createElement('div');
    hotelContainer.className = 'hotel';

    const hotelName = document.createElement('h2');
    hotelName.textContent = hotel.name;

    const hotelAddress = document.createElement('p');
    hotelAddress.textContent = `Country: ${hotel.address.countryCode}`;

    const hotelDistance = document.createElement('p');
    hotelDistance.textContent = `Distance: ${hotel.distance.value} ${hotel.distance.unit}`;

    hotelContainer.appendChild(hotelName);
    hotelContainer.appendChild(hotelAddress);
    hotelContainer.appendChild(hotelDistance);
    hotelInfoContainer.appendChild(hotelContainer);
  });
}

document.addEventListener("DOMContentLoaded", function () {
  fetch('/hotel')
    .then(response => response.json())
    .then(data => {
      const hotels = data.data.slice(0, 10);
      updateHotelInfo(hotels); 
    })
    .catch(error => {
      console.error('Error fetching hotel data:', error);
    });
});

document.addEventListener("DOMContentLoaded", function () {
  fetch('/activitie')
    .then(response => response.json())
    .then(data => {
      const activitie = data.data.slice(0, 5);
      updateActivitieInfo(activitie); 
    })
    .catch(error => {
      console.error('Error fetching activitie data:', error);
    });
});
