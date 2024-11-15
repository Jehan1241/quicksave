import { useState, useEffect } from 'react'
import ListFoundGames from '../AddGameManually/ListFoundGames'

function SteamImportView(props) {
  const [searchClicked, setSearchClicked] = useState(false)
  const [data, setData] = useState('')
  const [toSearch, setToSearch] = useState('')
  const [loading, setLoading] = useState(false)
  const [clickCount, setClickCount] = useState(0)
  const [selectedPlatform, setSelectedPlatform] = useState('')
  const [timePlayed, setTimePlayed] = useState(0)
  const [platforms, setPlatforms] = useState([])
  const [customPlatfromClicked, setCustomPlatformClicked] = useState(false)
  const [clientID, setClientID] = useState('')
  const [clientSecret, setClientSecret] = useState('')

  useEffect(() => {
    if (searchClicked) {
      fetchData(toSearch)
    }
  }, [clickCount])

  useEffect(() => {
    if (data !== '') {
      console.log(Object.values(data))
    }
  }, [data])

  const searchClickHandler = () => {
    setClickCount(clickCount + 1)
    setSearchClicked(true)
    const value = document.getElementById('SearchBar').value
    const time = document.getElementById('timePlayed').value
    const platform = document.getElementById('Platform').value
    const ClientID = document.getElementById('ClientID').value
    const ClientSecret = document.getElementById('ClientSecret').value
    setClientID(ClientID)
    setClientSecret(ClientSecret)
    setSelectedPlatform(platform)
    setTimePlayed(time)
    setLoading(true)
    setToSearch(value)
  }

  const handleTimePlayedChange = (e) => {
    const { value } = e.target
    const numericValue = value.replace(/[^0-9]/g, '')
    setTimePlayed(numericValue)
  }

  const fetchData = async (toSearch) => {
    try {
      const response = await fetch(`http://localhost:8080/IGDBsearch`, {
        method: 'POST',
        headers: { 'Content-type': 'application/json' },
        body: JSON.stringify({
          NameToSearch: toSearch,
          clientID: clientID,
          clientSecret: clientSecret
        })
      })
      setData(await response.json())
      setLoading(false)
    } catch (error) {
      console.error(error)
    }
  }

  useEffect(() => {
    getPlatforms()
  }, [])

  const getPlatforms = async () => {
    console.log('Getting Platforms')
    try {
      const response = await fetch(`http://localhost:8080/Platforms`)
      const json = await response.json()
      setPlatforms(Object.values(json.platforms))
      console.log('Run')
    } catch (error) {
      console.error(error)
    }
  }

  const checkForEnterPressed = (e) => {
    if (e.key == 'Enter') {
      searchClickHandler()
    }
  }

  const addCustomPlatformClickHandler = () => {
    setCustomPlatformClicked(!customPlatfromClicked)
  }

  return (
    <div className="flex flex-col gap-4 p-6 mt-4 w-full h-full text-base rounded-xl">
      <div className="flex flex-row gap-6 w-full h-full">
        <div className="flex flex-col gap-4 w-1/2">
          <div className="flex flex-row gap-4 items-center">
            <p className="w-24">Title</p>
            <input
              onKeyDown={checkForEnterPressed}
              id="SearchBar"
              className="px-2 w-52 rounded-lg bg-gray-500/20"
            />
          </div>

          <div className="flex flex-row gap-4 items-center">
            <p className="w-24">Hours Played</p>
            <input
              id="timePlayed"
              type="text"
              value={timePlayed}
              onChange={handleTimePlayedChange}
              className="px-2 w-52 h-6 rounded-lg bg-gray-500/20"
              onKeyDown={checkForEnterPressed}
            />
          </div>

          <div className="flex flex-row gap-4 items-center">
            <p className="w-24">Platform</p>
            {customPlatfromClicked ? (
              <input
                onKeyDown={checkForEnterPressed}
                id="Platform"
                className="px-2 w-52 h-6 rounded-lg bg-gray-500/20"
              />
            ) : (
              <select
                onKeyDown={checkForEnterPressed}
                id="Platform"
                className="px-2 w-52 h-6 rounded-lg bg-gray-500/20"
              >
                {platforms.map((item) => (
                  <option key={item.id} value={item.name}>
                    {item}
                  </option>
                ))}
              </select>
            )}
            <button
              className="px-2 h-6 text-sm rounded-lg bg-gray-500/20"
              onClick={addCustomPlatformClickHandler}
            >
              +
            </button>
          </div>

          <div className="flex flex-row gap-4 items-center">
            <p className="w-24">ClientID</p>
            <input
              id="ClientID"
              type="text"
              className="px-2 w-52 h-6 rounded-lg bg-gray-500/20"
              onKeyDown={checkForEnterPressed}
            />
          </div>

          <div className="flex flex-row gap-4 items-center">
            <p className="w-24">Client Secret</p>
            <input
              id="ClientSecret"
              type="text"
              className="px-2 w-52 h-6 rounded-lg bg-gray-500/20"
              onKeyDown={checkForEnterPressed}
            />
            <a href="https://api-docs.igdb.com/#getting-started" className="ml-2 text-sm">
              Client?
            </a>
          </div>
        </div>

        <div className="ml-auto w-full rounded-lg border border-white/10">
          {loading ? (
            <div className="text-left bg-gameView h-[28vh] w-auto overflow-scroll flex justify-center items-center">
              <p>Loading...</p>
            </div>
          ) : (
            <ListFoundGames
              FoundGames={data}
              SelectedPlatform={selectedPlatform}
              onGameAdded={props.onGameAdded}
              time={timePlayed}
            />
          )}
        </div>
      </div>

      <div className="flex justify-end mt-4">
        <button
          className="w-32 h-10 rounded-lg border-2 bg-gameView hover:bg-gray-500/20"
          onClick={searchClickHandler}
        >
          Search
        </button>
      </div>
    </div>
  )
}

export default SteamImportView
