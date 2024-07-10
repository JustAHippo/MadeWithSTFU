import './style.css';
import './app.css';

import logo from './assets/images/logo-universal.png';
import {FolderSelector, RunDeletion} from '../wailsjs/go/main/App';

document.querySelector('#app').innerHTML = `
    <img id="logo" class="logo">
      <div class="result">Select the Unity project's folder</div>
      
      <div class="input-box" id="input" style="padding-bottom: 25px;">
        Folder: <input class="input" id="name" type="text" autocomplete="off" /> 
        <button class="btn" onclick="folderSelector()">Select</button>
        
      </div>
      
      <div class="input-box"> 
        Version: <input class="input" id="version" type="text" autocomplete="off" />
        
        <button class="btn" onclick="runDeletion()">Go!</button>
      </div>
      <div class="result" id="result"></div>
    </div>
`;
document.getElementById('logo').src = logo;

let nameElement = document.getElementById("name");
nameElement.focus();
let resultElement = document.getElementById("result");

// Opens the folder selector
window.folderSelector = function () {

    try {
        FolderSelector(name)
            .then((result) => {

                nameElement.value = result;
            })
            .catch((err) => {
                console.error(err);
            });
    } catch (err) {
        console.error(err);
    }
};

window.runDeletion = function () {


    name = document.getElementById("name").value;
    let version;
    version = document.getElementById("version").value;

    try {
        RunDeletion(name, version)
            .then((result) => {
                resultElement.innerText = result;
            })
            .catch((err) => {
                console.error(err);
            });
    } catch (err) {
        console.error(err);
    }
};