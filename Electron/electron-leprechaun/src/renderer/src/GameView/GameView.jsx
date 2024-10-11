import DisplayInfo from './DisplayInfo'
import DisplayImage from './DisplayImage'
import { MdDelete } from 'react-icons/md'
import { useEffect, useState } from 'react'
import { FaPlay } from 'react-icons/fa'
import { TiTick } from 'react-icons/ti'
import { RxCross2 } from 'react-icons/rx'
import { useNavigate } from 'react-router-dom'
import AddGamePath from './AddGamePath'

function GameView(props) {
  const navigate = useNavigate()
  const [companies, setCompanies] = useState('')
  const [tags, setTags] = useState('')
  const [screenshots, setScreenshots] = useState('')
  const [metadata, setMetadata] = useState('')
  const [deleteClicked, setDeleteClicked] = useState(false)
  const [toAddGamePath, setToAddGamePath] = useState(false)

  const deleteGameClickHandler = () => {
    setDeleteClicked(!deleteClicked)
  }

  const confirmDeleteClickHandler = async () => {
    console.log('ABC')
    try {
      const response = await fetch(`http://localhost:8080/DeleteGame?uid=${props.uid}`)
      const json = await response.json()
    } catch (error) {
      console.error(error)
    }
    navigate('/', { replace: true })
    props.onDelete()
  }

  const addGamePathClickHandler = () => {
    console.log('here')
    setToAddGamePath(!toAddGamePath)
  }

  const sendGamePathtoDB = async (path) => {
    addGamePathClickHandler()
    console.log(path)
    try {
      const response = await fetch(
        `http://localhost:8080/setGamePath?uid=${props.uid}&path=${path}`
      )
      const json = await response.json()
    } catch (error) {
      console.error(error)
    }
  }

  const fetchData = async () => {
    try {
      const response = await fetch(`http://localhost:8080/GameDetails?uid=${props.uid}`)
      const json = await response.json()
      // Destructure metadata from the JSON response
      const { companies, tags, screenshots, m: metadata } = json.metadata
      // Set state correctly
      setCompanies(companies[props.uid]) // Access companies by UID
      setTags(tags[props.uid]) // Access tags by UID
      setMetadata(metadata[props.uid]) // Access metadata by UID
      setScreenshots(screenshots[props.uid]) // Access screenshots by UID
    } catch (error) {
      console.error(error)
    }
  }

  useEffect(() => {
    fetchData()
  }, [])

  const playClicked = async () => {
    try {
      const response = await fetch(`http://localhost:8080/LaunchGame?uid=${props.uid}`)
      const json = await response.json()
      console.log(json.ManualGameLaunch)
      if (json.ManualGameLaunch == 'AddPath') {
        addGamePathClickHandler()
      }
    } catch (error) {
      console.log(error)
    }
  }

  const UID = props.uid
  const tagsArray = Object.values(tags)
  const companiesArray = Object.values(companies)
  const screenshotsArray = Object.values(screenshots)

  console.log('Time', metadata.TimePlayed)

  return (
    <>
      <img
        className="absolute top-0 right-0 w-screen h-screen opacity-20 blur-md"
        src={'http://localhost:8080/screenshots/' + screenshots[0]}
      />
      <div className="overflow-y-auto relative h-screen text-center text-white">
        {/* Spacer Div */}
        <div className="h-[4%]"></div>
        {/* Name and Time Flex */}
        <div className="flex flex-row justify-between items-center px-3 mx-10 my-2 h-24 rounded-2xl">
          <div className="flex flex-row gap-3">
            <div className="text-3xl font-bold">{metadata.Name}</div>
            <button className="px-1 mt-1" onClick={playClicked}>
              <FaPlay size={18} />
            </button>
            <div className="flex flex-row items-center">
              <button className="px-1 mt-1" onClick={deleteGameClickHandler}>
                <MdDelete size={22} />
              </button>
              {deleteClicked ? (
                <div className="flex flex-row gap-2 mt-2 text-sm">
                  Are you sure?
                  <button onClick={confirmDeleteClickHandler}>
                    <TiTick size={20} />
                  </button>
                  <button onClick={deleteGameClickHandler}>
                    <RxCross2 size={20} />
                  </button>
                </div>
              ) : null}
            </div>
          </div>
          <div className="text-2xl">Time Played : {metadata.TimePlayed} Hrs</div>
        </div>
        {/* Horizontal FLEX Holder */}
        <div className="flex flex-row mx-10">
          {/* Description DIV */}
          <div className="m-2 w-1/3 h-full rounded-3xl">
            <DisplayInfo data={metadata} tags={tagsArray} companies={companiesArray} />
          </div>
          {/* Image Div */}
          <DisplayImage screenshots={screenshotsArray} />
        </div>
      </div>
      {toAddGamePath ? (
        <AddGamePath
          addGamePathClickHandler={addGamePathClickHandler}
          sendGamePathtoDB={sendGamePathtoDB}
        />
      ) : null}
    </>
  )
}

export default GameView
