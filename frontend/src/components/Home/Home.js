import React from 'react';
import Form from "@rjsf/material-ui";
import {baseSchema, uiSchema, customFormats} from './schemas.js';
import { useHistory } from 'react-router-dom';

var functions = require('../functions');

const Home = () => {
    const history = useHistory();

    const submitData = ({ formData }, e) => {
        functions.sendMsg(JSON.stringify(formData));
        history.push("/deploy");
    }

    return (
    <div className="App">
        <Form schema={baseSchema} onSubmit={submitData} customFormats={customFormats} uiSchema={uiSchema} />
    </div>
    );
}
 
export default Home;