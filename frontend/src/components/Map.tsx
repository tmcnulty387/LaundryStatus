import { Link } from "react-router-dom";
import mapImage from "../assets/map.png";

function Map() {
  return (
    <main className="home">
      <div className="map-wrapper">
        <img
          src={mapImage}
          alt="Map of Residence Halls with Waypoints"
          className="campus-map"
        />
        <Link
          to="/residence-hall-a"
          className="waypoint"
          style={{ top: "39.1%", left: "25.2%" }}
        >
          <span className="waypoint-dot"></span>
          <span className="waypoint-label">Residence Hall A</span>
        </Link>
        <Link
          to="/residence-hall-b"
          className="waypoint"
          style={{ top: "39.1%", left: "32.1%" }}
        >
          <span className="waypoint-dot"></span>
          <span className="waypoint-label">Residence Hall B</span>
        </Link>
        <Link
          to="/residence-hall-c"
          className="waypoint"
          style={{ top: "41.9%", left: "44.9%" }}
        >
          <span className="waypoint-dot"></span>
          <span className="waypoint-label">Residence Hall C</span>
        </Link>
        <Link
          to="/kate-gleason"
          className="waypoint"
          style={{ top: "61.6%", left: "22.6%" }}
        >
          <span className="waypoint-dot"></span>
          <span className="waypoint-label">Kate Gleason</span>
        </Link>
        <Link
          to="/gibson"
          className="waypoint"
          style={{ top: "64.9%", left: "55.3%" }}
        >
          <span className="waypoint-dot"></span>
          <span className="waypoint-label">Gibson</span>
        </Link>
        <Link
          to="/peterson"
          className="waypoint"
          style={{ top: "63.2%", left: "63.1%" }}
        >
          <span className="waypoint-dot"></span>
          <span className="waypoint-label">Peterson</span>
        </Link>
        <Link
          to="/sol-heumann"
          className="waypoint"
          style={{ top: "77.0%", left: "44.6%" }}
        >
          <span className="waypoint-dot"></span>
          <span className="waypoint-label">Sol Heumann</span>
        </Link>
      </div>
    </main>
  );
}

export default Map;
