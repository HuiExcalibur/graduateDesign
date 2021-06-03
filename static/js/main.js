// import axios from 'axios';
// import Vue from "vue";

Vue.prototype.$axios = axios;

var sock=null
var wsURL="ws://127.0.0.1:8080/WS"

var app = new Vue({
	delimiters:['${','}'],
	el:"#vueapp",
	data:function(){
		return{
			username:"",
			nickname:"",
			input:"",
			rooms:[],
			room_name:"",
			search_result:[],
			NDVisibility:false,
			new_roomname:"",
			SDVisibility:false,
			new_nickname:"",
			CDVisibility:false,
			search_keyword:"",
			messages:[],
			// messages:[{
			// 	message:"nihao?????????????????????????????????????????????????????????????????????????????????????????",
			// 	username:"other",
			// 	roomname:"room1"
			// },{
			// 	message:"woyehao",
			// 	username:"me",
			// 	roomname:"room1"
			// }],
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
				alert("socket is closed, please try again")
				this.init_websocket()
			}
			// this.messages.push({
			// 	msg:this.input,
			// 	user:this.username
			// })
		},
		create_room:function(){
			var that=this
			var url='/newroom?roomname='+that.new_roomname+'&username='+that.username

			axios.get(url)
			.then(function(response){
				console.log(response.data)
				if (response.data.status=='success'){
					that.rooms.push(that.new_roomname)
					that.$message({
						message:'create success',
						type:'success',
					})
				}else{
					that.$message({
						message:'failure',
						type:'error',
					})
				}
			})
			.catch(function(err){
				console.log(err)
			})

			this.NDVisibility=false
		},
		quit_room:function(roomname){
			var that=this
			var url='/quit?roomname='+roomname+'&username='+that.username

			axios.get(url)
			.then(function(response){
				console.log(response.data)
				if (response.data.status=='success'){
					that.$message({
						message:'quit success',
						type:'success',
					})
					for (var i=0;i<that.rooms.length;i++){
						if (that.rooms[i]==roomname){
							that.rooms.splice(i,1)
						}
					}
				}else{
					that.$message({
						message:'failure',
						type:'error',
					})
				}
			})
			.catch(function(err){
				console.log(err)
			})
		},
		search_room:function(){
			var that=this
			var url='/search?key='+that.search_keyword

			axios.get(url)
			.then(function(response){
				console.log(response.data)
				that.search_result=response.data.rooms
			})
			.catch(function(err){
				console.log(err)
			})

			that.SDVisibility=false
		},
		enter_room:function(roomindex){
			var that=this
			var url='/enter?roomname='+that.search_result[roomindex]
			console.log(url)
			
			axios.get(url)
			.then(function(response){
				if (response.data.status=='success'){
					that.$message({
						message:'加入成功',
						type:'success',
					})
					that.rooms.push(that.search_result[roomindex])
					that.get_history(that.search_result[roomindex])
					that.search_result.splice(roomindex,1)
				}else{
					that.$message({
						message:'failure',
						type:'error',
					})
				}
			})
			.catch(function(err){
				console.log(err)
			})
		},
		get_history:function(roomname){
			var url='/history?roomname='+roomname

			axios.get(url)
			.then(function(response){
				console.log(response.data)
			})
			.catch(function(err){
				console.log(err)
			});
		},
		change_nickname:function(){
			var url='/changenickname?nickname='+this.new_nickname

			var that=this
			axios.get(url)
			.then(function(response){
				console.log(response.data)
				if (response.data.status=="success"){
					that.nickname=that.new_nickname
				}
			})
			.catch(function(err){
				console.log(err)
			})

			that.CDVisibility=false;
		},
		init_data:function(){
			var cn=document.cookie.split(';')
            for (var i=0;i<cn.length;i++){
                value=cn[i].trim().split('=')
                        
                console.log("key ",value[0],"value",value[1])
				if (value[0]=='user'){
					this.username= decodeURI(value[1])
				}
				if (value[0]=='nickname'){
					this.nickname= decodeURI(value[1])
				}
            }

			var that=this
			that.$axios.get('/getroom')
            .then(function(response){
                console.log(response.data)
                // that.rooms=that.rooms.concat(response.data.rooms)
				that.rooms=response.data.rooms
				that.room_name=that.rooms[0]
				for (var i=0;i<that.rooms.length;i++){
					that.get_history(that.rooms[i])
				}
            })
            .catch(function(err){
				console.log(err)
            });
		},
		init_websocket:function(){
			var that=this

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
		}
	},


	updated : function(){
		this.$nextTick(function(){
			var div = document.getElementById('show_messages');
			div.scrollTop = div.scrollHeight;
		})

	},
	created:function(){
		// var that=this
		console.log("Vue is created",this.test)
		this.init_websocket()
		this.init_data()
	},
	destroyed:function(){
		sock.close()
	}
})