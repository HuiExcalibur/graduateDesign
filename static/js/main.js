

var sock=null
var wsURL="ws://127.0.0.1:8080/WS"

var app = new Vue({
	delimiters:['${','}'],
	el:"#vueapp",
	data:function(){
		return{
			username:"me",
			input:"",
			rooms:["room1","room2","room3"],
			room_name:"mainRoom",
			messages:[{
				message:"nihao?????????????????????????????????????????????????????????????????????????????????????????",
				username:"other",
				roomname:"room1"
			},{
				message:"woyehao",
				username:"me",
				roomname:"room1"
			}],
		}
	},
	methods:{
		read:function(room){
			console.log("now i'm in room",room)
			this.room_name=room
		},
		send_message:function(){
			//this.rooms.push(this.input)
			var msg_send={
				username:this.username,
				message:this.input,
				roomname:this.room_name
			}
			if (sock.readyState==1){
				sock.send(JSON.stringify(msg_send))
				this.input=""
			}else{
				alert("socket is closed")
			}
			// this.messages.push({
			// 	msg:this.input,
			// 	user:this.username
			// })
		}
	},
	updated : function(){
		this.$nextTick(function(){
			var div = document.getElementById('show_messages');
			div.scrollTop = div.scrollHeight;
		})

	},
	created:function(){
		var that=this
		console.log("Vue is created",that.test)
		if(typeof(WebSocket)==="undefined"){
			alert("不支持websocket")
			return
		}
		sock=new WebSocket(wsURL)
		sock.onopen=function(){
			console.log("connect to "+wsURL)
		}
		sock.onclose=function(e){
			console.log("close the sock with code "+e.code)
		}
		sock.onmessage=function(e){
			console.log("receive message "+e.data)
			receive=JSON.parse(e.data)
			// if (receive.username==that.input.username){
			// 	receive.username=that.username
			// }
			that.messages.push(receive)
		}
	},
	destroyed:function(){
		sock.close()
	}
})