import { useState } from 'react'
import MetaData from './MetaData'
import Delete from './Delete'

function Customize(props) {
  const [metaDataClicked, setMetaDataClicked] = useState(true)
  const [imagesClicked, setImagesClicked] = useState(false)
  const [deleteClicked, setDeleteClicked] = useState(false)

  const deleteClickHandler = () => {
    setDeleteClicked(true)
    setImagesClicked(false)
    setMetaDataClicked(false)
  }

  const metaDataClickHandler = () => {
    setMetaDataClicked(true)
    setDeleteClicked(false)
    setImagesClicked(false)
  }

  return (
    <div
      className="flex overflow-hidden fixed top-0 right-0 w-screen h-screen text-white bg-black/80"
      onClick={props.customizeClickHandler}
    >
      <div
        className="flex flex-col p-5 m-auto w-2/3 h-2/3 text-2xl bg-gameView"
        onClick={(e) => e.stopPropagation()}
      >
        <div className="pb-2 mb-2 border-b-2">Customize Game Data</div>
        <div className="flex flex-row">
          <button
            onClick={metaDataClickHandler}
            className={`mx-2 h-10 text-lg ${metaDataClicked ? 'border-b-2' : null}`}
          >
            Metadata
          </button>

          <button className={`mx-2 h-10 text-lg ${imagesClicked ? 'border-b-2' : null}`}>
            Images
          </button>

          <button
            onClick={deleteClickHandler}
            className={`mx-2 h-10 text-lg ${imagesClicked ? 'border-b-2' : null}`}
          >
            Delete
          </button>
        </div>
        {metaDataClicked && <MetaData uid={props.uid} re_renderTrigger={props.re_renderTrigger} />}
        {deleteClicked && <Delete uid={props.uid} />}
      </div>
    </div>
  )
}

export default Customize
