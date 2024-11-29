
import { useRef, useState, useCallback, useEffect } from 'react'
import { DeBoxChatWidget } from '@debox-pro/chat-widget-react' // Import DeBoxChatWidget component

export default function App() {
  const inputRef = useRef(null)
  const [cid, setCID] = useState('')
  const [msg, setMsg] = useState([])
  const [group_name, setGroupName] = useState('op chain usdt')
  const [chain_id, setChainId] = useState(10)
  const [contract_address, setContractAddress] = useState('0xb9a96adc79d4d9d3c7f3455ff5b2a98b5915b0a5')

  // Get query parameters from URL
  const getQueryVariable = useCallback((variable) => {
    return new URL(window.location.href).searchParams.get(variable)
  }, [])

  // Initialize conversation ID
  const initializeConversationId = useCallback(async () => {
    const groupName = getQueryVariable('group_ame')
    const chainId = getQueryVariable('chain_id')
    const contractAddress = getQueryVariable('contract_address')

    if (groupName || chainId || contractAddress) {
      try {
        const response = await fetch(
          `https://demo-chat.debox.pro/getconversationid?group_name=${encodeURIComponent(groupName || '')}&chain_id=${encodeURIComponent(chainId || '')}&contract_address=${encodeURIComponent(contractAddress || '')}`,
          {
            headers: {
              'Content-Type': 'application/json'
            }
          }
        )

        const data = await response.json()
        if (data.success) {
          setCID(data.data.cid)
        } else {
          console.warn("Failed to initialize conversation ID")
        }
      } catch (error) {
        console.error("Error initializing conversation ID:", error)
      }
    }
  }, [getQueryVariable])

  // Modify function to get conversation ID
  const fetchConversationId = useCallback(async () => {
    if (!group_name && !chain_id && !contract_address) {
      alert("Please fill in at least one field")
      return
    }

    try {
      const response = await fetch(
        `https://demo-chat.debox.pro/getconversationid?group_name=${encodeURIComponent(group_name)}&chain_id=${encodeURIComponent(chain_id)}&contract_address=${encodeURIComponent(contract_address)}`,
        {
          headers: {
            'Content-Type': 'application/json'
          }
        }
      )

      const data = await response.json()
      if (data.success) {
        const conversationId = data.data.cid
        if (inputRef.current) {
          inputRef.current.value = conversationId
        }
        alert(`Conversation ID: ${conversationId}`)
      } else {
        alert("Failed to get conversation ID")
      }
    } catch (error) {
      alert("Error getting conversation ID")
      console.error(error)
    }
  }, [group_name, chain_id, contract_address])

  // Initialize on component load
  useEffect(() => {
    initializeConversationId()
  }, [initializeConversationId])

  // Handle message events
  const handleEvent = useCallback((e) => {
    setMsg(prev => [...prev, e?.detail])
  }, [])


  // Update conversation ID
  const updateCID = useCallback(() => {
    setCID(inputRef.current?.value)
  }, [])

  // Add function to handle input changes
  const handleInputChange = useCallback((e) => {
    const { id, value } = e.target
    switch(id) {
      case 'group_name':
        setGroupName(value);
        break;
      case 'chain_id':
        setChainId(value);
        break;
      case 'contract_address':
        setContractAddress(value);
        break;
      default:
        break;
    }
  }, [])

  return (
    <div style={{ height: '100vh' }}>
      <div id="test-dom">
        <h2>This is the dApp page...</h2>
      </div>

      <h3>Get conversation ID based on input parameters:</h3>
      <div>
        <input type="text" id="group_name" placeholder="Enter group name" onChange={handleInputChange} />
        <input type="text" id="chain_id" placeholder="Enter chain ID" onChange={handleInputChange} />
        <input type="text" id="contract_address" placeholder="Enter contract address" onChange={handleInputChange} />
        <button onClick={fetchConversationId}>Get Conversation ID</button>
      </div>

      <h3>Input conversation ID to change chat:</h3>
      <div>
        <input
          ref={inputRef}
          type="text"
          id="conversationInput"
          placeholder="Enter conversation ID here"
        />
        <button onClick={updateCID}>Change Conversation</button>
      </div>



        <DeBoxChatWidget
          projectId={process.env.REACT_APP_PKEY}
          conversationId={cid}
          onEvent={handleEvent}
        />

      <h3>Event Log:</h3>
      <div>
        {msg.map((m, i) => (
          <li key={i}>{JSON.stringify(m)}</li>
        ))}
      </div>
    </div>
  )
}
