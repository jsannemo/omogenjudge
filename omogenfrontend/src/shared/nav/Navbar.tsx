function Navbar() {
    return (
        <nav className="navbar navbar-dark bg-dark navbar-expand-lg" id="menu">
            <div className="container">
                <a className="navbar-brand me-5" href="/">
                    <img src="/static/img/logo.png" height="45" alt="Logotype" />
                </a>
                <button className="navbar-toggler" type="button" data-bs-toggle="collapse"
                        data-bs-target="#navbarSupportedContent" aria-controls="navbarSupportedContent"
                        aria-expanded="false" aria-label="Toggle navigation">
                    <span className="navbar-toggler-icon"/>
                </button>
                <div className="collapse navbar-collapse" id="navbarSupportedContent">
                    <ul className="navbar-nav me-auto">
                        <a className="nav-link " href="/submissions/">Submissions</a>
                    </ul>
                    <ul className="navbar-nav">
                        <a className="nav-link " href="/accounts/logout/">Log out</a>
                    </ul>
                </div>
            </div>
        </nav>
    );
}

export default Navbar;