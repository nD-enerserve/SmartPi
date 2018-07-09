/*
    Copyright (C) Jens Ramhorst
  	This file is part of SmartPi.
    SmartPi is free software: you can redistribute it and/or modify
    it under the terms of the GNU General Public License as published by
    the Free Software Foundation, either version 3 of the License, or
    (at your option) any later version.
    SmartPi is distributed in the hope that it will be useful,
    but WITHOUT ANY WARRANTY; without even the implied warranty of
    MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
    GNU General Public License for more details.
    You should have received a copy of the GNU General Public License
    along with SmartPi.  If not, see <http://www.gnu.org/licenses/>.
    Diese Datei ist Teil von SmartPi.
    SmartPi ist Freie Software: Sie können es unter den Bedingungen
    der GNU General Public License, wie von der Free Software Foundation,
    Version 3 der Lizenz oder (nach Ihrer Wahl) jeder späteren
    veröffentlichten Version, weiterverbreiten und/oder modifizieren.
    SmartPi wird in der Hoffnung, dass es nützlich sein wird, aber
    OHNE JEDE GEWÄHRLEISTUNG, bereitgestellt; sogar ohne die implizite
    Gewährleistung der MARKTFÄHIGKEIT oder EIGNUNG FÜR EINEN BESTIMMTEN ZWECK.
    Siehe die GNU General Public License für weitere Details.
    Sie sollten eine Kopie der GNU General Public License zusammen mit diesem
    Programm erhalten haben. Wenn nicht, siehe <http://www.gnu.org/licenses/>.
*/
/*
File: measurement.go
Description: Defines the structure of of the json-api-file
*/

package smartpi

type tValue struct {
	Type  string  `json:"type" xml:"type"`
	Unity string  `json:"unity" xml:"unity"`
	Info  string  `json:"info" xml:"info"`
	Data  float32 `json:"data" xml:"data"`
}

type tPhase struct {
	Phase  int       `json:"phase" xml:"phase"`
	Name   string    `json:"name" xml:"name"`
	Values []*tValue `json:"values" xml:"values"`
}

type tDataset struct {
	Time   string    `json:"time" xml:"time"`
	Phases []*tPhase `json:"phases" xml:"phases"`
}

type tMeasurement struct {
	Serial          string      `json:"serial" xml:"serial"`
	Name            string      `json:"name" xml:"name"`
	Lat             float64     `json:"lat" xml:"lat"`
	Lng             float64     `json:"lng" xml:"lng"`
	Time            string      `json:"time" xml:"time"`
	Softwareversion string      `json:"softwareversion" xml:"softwareversion"`
	Ipaddress       string      `json:"ipaddress" xml:"ipaddress"`
	Datasets        []*tDataset `json:"datasets" xml:"datasets"`
}
