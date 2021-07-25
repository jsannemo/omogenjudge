import React from "react";
import {Container, Nav, Navbar} from "react-bootstrap";
import {Link} from "react-router-dom";

function Menu() {
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
          <Nav>
            <Nav.Link href={"/login"}>Login</Nav.Link>
            <Nav.Link as={Link} to={"/register"}>Register</Nav.Link>
          </Nav>
        </Navbar.Collapse>
      </Container>
    </Navbar>
  );
}


export default Menu;