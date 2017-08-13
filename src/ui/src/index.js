import React from 'react';
import ReactDOM from 'react-dom';
import './index.css';

class UploadForm extends React.Component {
	render() {
		return (
				<form action="http://localhost:8080/api/v1/import/fitnotes" method="post" encType="multipart/form-data">
				<input type="file" name="files" multiple />
				<input type="submit" value="Submit" />
				</form>
		       );
	}
}

ReactDOM.render(
		<UploadForm / >,
		document.getElementById('root')
	       );
