import { useEffect, useState } from 'react'

function MetaData(props) {
  const [customTitleChecked, setCustomTitleChecked] = useState(false)
  const [customTitle, setCustomTitle] = useState('')
  const [customTime, setCustomTime] = useState('')
  const [customTimeOffset, setCustomTimeOffset] = useState('')
  const [customTimeChecked, setCustomTimeChecked] = useState(false)
  const [customTimeOffsetChecked, setCustomTimeOffsetChecked] = useState(false)

  const loadPreferences = async () => {
    try {
      const response = await fetch(`http://localhost:8080/LoadPreferences?uid=${props.uid}`)
      const json = await response.json()
      console.log(json)
      if (json.preferences.title.checked == '1') {
        setCustomTitleChecked(true)
        setCustomTitle(json.preferences.title.value)
      }
      if (json.preferences.time.checked == '1') {
        setCustomTimeChecked(true)
        setCustomTime(json.preferences.time.value)
      }
      if (json.preferences.timeOffset.checked == '1') {
        setCustomTimeOffsetChecked(true)
        setCustomTimeOffset(json.preferences.timeOffset.value)
      }
      console.log(json.preferences.title.value)
    } catch (error) {
      console.error(error)
    }
  }

  useEffect(() => {
    loadPreferences()
  }, [])

  const saveClickHandler = () => {
    const customTitleElement = document.getElementById('customTitle')
    let title = ''
    let time = '0'
    let timeOffset = '0'
    const uid = props.uid

    if (customTitleElement && customTitleElement.value.trim() !== '') {
      title = customTitleElement.value
    }

    if (customTimeChecked) {
      time = document.getElementById('customTime').value.toString()
    }
    if (customTimeOffsetChecked) {
      timeOffset = document.getElementById('customTimeOffset').value.toString()
    }

    console.log(title, time, timeOffset)

    const postData = {
      customTitleChecked: customTitleChecked,
      customTitle: title,
      customTimeChecked: customTimeChecked,
      customTime: time,
      customTimeOffsetChecked: customTimeOffsetChecked,
      customTimeOffset: timeOffset,
      UID: uid
    }
    savePreferences(postData)
  }

  const savePreferences = async (postData) => {
    try {
      const response = await fetch(`http://localhost:8080/SavePreferences`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json'
        },
        body: JSON.stringify(postData)
      })
      const json = await response.json()
      if (json.status === 'OK') {
        props.re_renderTrigger()
      }
    } catch (error) {
      console.error(error)
    }
  }

  const handleCheckboxChange = (event, title) => {
    const checked = event.target.checked
    switch (title) {
      case 'title':
        setCustomTitleChecked(checked)
        break

      case 'time':
        if (checked == false) {
          setCustomTimeChecked(false)
        } else {
          setCustomTimeChecked(true)
          setCustomTimeOffsetChecked(false)
        }
        break

      case 'timeOffset':
        if (checked == false) {
          setCustomTimeOffsetChecked(false)
        } else {
          setCustomTimeOffsetChecked(true)
          setCustomTimeChecked(false)
        }
        break

      default:
        break
    }
  }

  return (
    <div className="flex flex-col gap-4 p-6 mt-4 w-full h-full text-base rounded-xl">
      <div className="flex flex-row gap-6 w-full h-full">
        <div className="flex flex-col gap-4 w-1/2">
          <div className="flex flex-row gap-4 items-center">
            <p>Use Custom Title?</p>
            <input
              type="checkbox"
              checked={customTitleChecked}
              onChange={(event) => handleCheckboxChange(event, 'title')}
              className="px-2 rounded-lg bg-gray-500/20"
            />
            {customTitleChecked && (
              <input
                type="text"
                id="customTitle"
                className="px-2 w-52 h-6 rounded-lg bg-gray-500/20"
                value={customTitle}
                onChange={(e) => setCustomTitle(e.target.value)}
              />
            )}
          </div>

          <div className="flex flex-row gap-4 items-center">
            <p>Set Custom Time Played?</p>
            <input
              type="checkbox"
              checked={customTimeChecked}
              onChange={(event) => handleCheckboxChange(event, 'time')}
              className="px-2 rounded-lg bg-gray-500/20"
            />
            {customTimeChecked && (
              <input
                type="number"
                id="customTime"
                className="px-2 w-52 h-6 rounded-lg bg-gray-500/20"
                value={customTime}
                onChange={(e) => setCustomTime(e.target.value)}
              />
            )}
          </div>

          <div className="flex flex-row gap-4 items-center">
            <p>Set Custom Time Offset?</p>
            <input
              type="checkbox"
              checked={customTimeOffsetChecked}
              onChange={(event) => handleCheckboxChange(event, 'timeOffset')}
              className="px-2 rounded-lg bg-gray-500/20"
            />
            {customTimeOffsetChecked && (
              <input
                type="number"
                id="customTimeOffset"
                className="px-2 w-52 h-6 rounded-lg bg-gray-500/20"
                value={customTimeOffset}
                onChange={(e) => setCustomTimeOffset(e.target.value)}
              />
            )}
          </div>
        </div>
      </div>

      <div className="flex justify-end mt-4">
        <button
          onClick={saveClickHandler}
          className="w-32 h-10 rounded-lg border-2 bg-gameView hover:bg-gray-500/20"
        >
          Save
        </button>
      </div>
    </div>
  )
}

export default MetaData
