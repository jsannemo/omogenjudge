import React from "react";
import {Container, Nav, Navbar} from "react-bootstrap";
import {Link} from "react-router-dom";
import {AuthStore} from "../../accounts/auth";

export default function Menu(): JSX.Element {
  const authState = AuthStore.useState();
  return (
    <Navbar
      variant={"dark"}
      bg={"dark"}
      expand={"lg"}>
      <Container>
        <Navbar.Brand className={"me-5"}>
          <img src="/static/img/logo.png" height="45"/>
        </Navbar.Brand>
        <Navbar.Toggle aria-controls="basic-navbar-nav"/>
        <Navbar.Collapse id="basic-navbar-nav">
          <Nav className="me-auto">
            <Nav.Link href="#home">Problems</Nav.Link>
          </Nav>
          {
            authState.profile
              ? <Nav>
                <Nav.Link as={Link} to={"/user"}>
                  {authState.profile.username}
                </Nav.Link>
                <Nav.Link as={Link} to={"/logout"}>Logout</Nav.Link>
              </Nav>
              : <Nav>
                <Nav.Link as={Link} to={"/login"}>Login</Nav.Link>
                <Nav.Link as={Link} to={"/register"}>Register</Nav.Link>
              </Nav>
          }
        </Navbar.Collapse>
      </Container>
    </Navbar>
  );
}