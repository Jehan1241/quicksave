import { useState, useEffect } from 'react'

function PsImportView() {
  const [clientID, setClientID] = useState('')
  const [clientSecret, setClientSecret] = useState('')
  const [npsso, setNpsso] = useState('')

  const searchClickHandler = async () => {
    const npsso = document.getElementById('npsso').value
    const clientID = document.getElementById('clientID').value
    const clientSecret = document.getElementById('clientSecret').value
    console.log(npsso)

    try {
      const response = await fetch('http://localhost:8080/PlayStationImport', {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify({ npsso: npsso, clientID: clientID, clientSecret: clientSecret })
      })
    } catch (error) {
      console.error('Error:', error)
    }
  }

  useEffect(() => {
    getIGDBKeys()
    getNpsso()
  }, [])

  const getIGDBKeys = async () => {
    console.log('Getting Platforms')
    try {
      const response = await fetch(`http://localhost:8080/IGDBKeys`)
      const json = await response.json()
      console.log(json.IGDBKeys[0])
      setClientID(json.IGDBKeys[0])
      setClientSecret(json.IGDBKeys[1])
    } catch (error) {
      console.error(error)
    }
  }

  const getNpsso = async () => {
    console.log('Getting Platforms')
    try {
      const response = await fetch(`http://localhost:8080/Npsso`)
      const json = await response.json()
      console.log(json.Npsso)
      setNpsso(json.Npsso)
    } catch (error) {
      console.error(error)
    }
  }

  const checkForEnterPressed = (e) => {
    if (e.key == 'Enter') {
      searchClickHandler()
    }
  }

  return (
    <div className="flex flex-col p-4 mt-2 w-full h-full text-base rounded-xl">
      {/* Input Div */}
      <div className="flex flex-col gap-4">
        <div className="flex flex-row gap-4">
          <div className="flex flex-col gap-4">
            <p>NPSSO</p>
          </div>
          <div className="flex flex-col gap-4">
            <input
              onKeyDown={checkForEnterPressed}
              id="npsso"
              className="px-1 mx-10 w-72 h-6 text-sm rounded-lg bg-gray-500/20"
              value={npsso}
              onChange={(e) => setNpsso(e.target.value)}
            ></input>
          </div>
          <div className="flex items-end text-sm text-blue-700 underline">
            <a className="text-sm" href="https://steamcommunity.com/dev/apikey">
              NPSSO?
            </a>
          </div>
        </div>
        <div className="flex flex-row gap-4">
          <div className="flex flex-col gap-4">
            <p>Client ID</p>
          </div>
          <div className="flex flex-col gap-4">
            <input
              onKeyDown={checkForEnterPressed}
              id="clientID"
              className="px-1 mx-8 w-72 h-6 text-sm rounded-lg bg-gray-500/20"
              value={clientID}
              onChange={(e) => setClientID(e.target.value)}
            ></input>
          </div>
        </div>
        <div className="flex flex-row gap-4">
          <div className="flex flex-col gap-4">
            <p>Client Secret</p>
          </div>
          <div className="flex flex-col gap-4">
            <input
              onKeyDown={checkForEnterPressed}
              id="clientSecret"
              className="px-1 w-72 h-6 text-sm rounded-lg bg-gray-500/20"
              value={clientSecret}
              onChange={(e) => setClientSecret(e.target.value)}
            ></input>
          </div>
        </div>
      </div>
      {/* button div */}
      <div className="flex justify-end mt-auto">
        <button className="w-32 h-10 rounded-lg border-2 bg-primary" onClick={searchClickHandler}>
          Import
        </button>
      </div>
    </div>
  )
}

export default PsImportView
