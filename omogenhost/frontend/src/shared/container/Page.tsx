import React from 'react';

import './Page.scss';

const Page: React.FunctionComponent = ({children}) =>
    (
        <div className="container page-content">
          <div className="mt-4">
            {children}
          </div>
        </div>
    );

export default Page;