import React from 'react'
import { Switch, Route } from 'react-router-dom'

import Home from './home/Home'
import Choose from './sets/Choose'
import Study from './sets/Study'

const Router = () => (
    <main>
        <Switch>
            <Route exact path='/' component={Home} />
            <Route path='/sets/choose' component={Choose} />
            <Route path='/sets/study' component={Study} />
        </Switch>
    </main>
)

export default Router;