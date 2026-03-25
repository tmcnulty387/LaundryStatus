import { Link } from "react-router-dom";

function Navbar() {
  return (
    <nav className="navbar">
      <Link className="brand-link" to="/">
        Laundry Status
      </Link>
      <div className="nav-room-links">
        <Link className="nav-room-link" to="/residence-hall-a">
          Residence Hall A
        </Link>
        <Link className="nav-room-link" to="/residence-hall-b">
          Residence Hall B
        </Link>
        <Link className="nav-room-link" to="/residence-hall-c">
          Residence Hall C
        </Link>
        <Link className="nav-room-link" to="/kate-gleason">
          Kate Gleason
        </Link>
        <Link className="nav-room-link" to="/gibson">
          Gibson
        </Link>
        <Link className="nav-room-link" to="/peterson">
          Peterson
        </Link>
        <Link className="nav-room-link" to="/sol-heumann">
          Sol Heumann
        </Link>
      </div>
    </nav>
  );
}

export default Navbar;
