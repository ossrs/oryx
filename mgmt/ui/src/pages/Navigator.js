import React from 'react';
import Container from "react-bootstrap/Container";
import {Navbar, Nav} from 'react-bootstrap';
import {Link, useLocation} from "react-router-dom";
import logo from '../resources/logo.svg';

export default function Navigator({initialized, token}) {
  const [activekey, setActiveKey] = React.useState(1);
  const [navs, setNavs] = React.useState([]);
  const location = useLocation();

  React.useEffect(() => {
    console.log(`xxx`, initialized, token, location.pathname);
    if (!initialized) return setNavs([]);

    if (!token) {
      return setNavs([{to:'/routers-login', text: '登录', className: 'text-light'}]);
    }

    setNavs([
      {eventKey: '1', to: '/routers-dashboard', text: '仪表盘'},
      {eventKey: '2', to: '/routers-scenario', text: '应用场景'},
      {eventKey: '3', to: '/routers-settings', text: '系统配置'},
      {eventKey: '4', to: '/routers-system', text: '组件管理'},
      {eventKey: '5', to: '/routers-contact', text: '专享群'},
    ].map(e => {
      if (e.to === location.pathname) {
        e.className = 'text-light';
        setActiveKey(e.eventKey);
      }
      return e;
    }));
  }, [initialized, token, location]);

  return (<>
    <Navbar>
      <Container className={{color:'#fff'}}>
        <Navbar.Brand>
          <img
            src={logo}
            width="64"
            height="30"
            className="d-inline-block align-top"
            alt="SRS Terraform"
          />
        </Navbar.Brand>
        <Nav className='me-auto' variant="pills" activeKey={activekey}>
          {navs.map((e, index) => {
            return (
              <Nav.Link
                as={Link}
                eventKey={e.eventKey}
                to={e.to}
                key={index}
                className={e.className}
              >
                {e.text}
              </Nav.Link>
            );
          })}
        </Nav>
      </Container>
    </Navbar>
  </>);
}

