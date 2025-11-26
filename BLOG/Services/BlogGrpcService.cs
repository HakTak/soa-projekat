using Blog; // This namespace comes from the Proto generation
using Grpc.Core;
using BLOG.Model;
using BLOG.Services;
using Google.Protobuf.WellKnownTypes;
using System.Linq; 

namespace BLOG.GrpcServices
{
    public class BlogGrpcService : BlogService.BlogServiceBase
    {
        private readonly PostService _postService;
        private readonly CommentService _commentService;

        public BlogGrpcService(PostService postService, CommentService commentService)
        {
            _postService = postService;
            _commentService = commentService;
        }

        // Helper to get UserID from Metadata (sent by Go Gateway)
        private string GetUserId(ServerCallContext context)
        {
            var userEntry = context.RequestHeaders.FirstOrDefault(h => h.Key == "user-id");
            
            if (userEntry == null || string.IsNullOrEmpty(userEntry.Value))
            {
                throw new RpcException(new Status(StatusCode.Unauthenticated, "No user ID found"));
            }
            return userEntry.Value;
        }

        public override async Task<GetAllPostsResponse> GetAllPosts(Empty request, ServerCallContext context)
        {
            var posts = await _postService.GetPostsAsync();
            var response = new GetAllPostsResponse();

            foreach (var post in posts)
            {
                response.Posts.Add(MapToProtoPost(post));
            }

            return response;
        }

        public override async Task<PostResponse> CreatePost(CreatePostRequest request, ServerCallContext context)
        {
            var newPost = new Post
            {
                Title = request.Title,
                Description = request.Description,
                ImagePaths = request.ImagePaths.ToArray(),
                CreatedAt = DateTime.UtcNow,
                LikeCount = 0
            };

            await _postService.CreatePostAsync(newPost);
            return MapToProtoPost(newPost);
        }

        public override async Task<PostResponse> ToggleLike(ToggleLikeRequest request, ServerCallContext context)
        {
            var userId = GetUserId(context); // Get from Header
            var updatedPost = await _postService.TogglePostLikeAsync(request.PostId, userId);

            if (updatedPost == null)
                throw new RpcException(new Status(StatusCode.NotFound, "Post not found"));

            return MapToProtoPost(updatedPost);
        }

        public override async Task<GetCommentsResponse> GetComments(GetCommentsRequest request, ServerCallContext context)
        {
            var comments = await _commentService.GetCommentsByPostIdAsync(request.PostId);
            var response = new GetCommentsResponse();

            foreach (var c in comments)
            {
                response.Comments.Add(MapToProtoComment(c));
            }
            return response;
        }

        public override async Task<CommentResponse> CreateComment(CreateCommentRequest request, ServerCallContext context)
        {
            var userId = GetUserId(context);
            // Optional: Get Username from header if you send it from Gateway, or fetch it
            var username = context.RequestHeaders.GetValue("x-user-username") ?? "Unknown"; 

            var newComment = new Comment
            {
                PostId = request.PostId,
                Text = request.Text,
                AuthorId = userId,
                AuthorName = username,
                CreatedAt = DateTime.UtcNow,
                UpdatedAt = DateTime.UtcNow
            };

            await _commentService.CreateCommentAsync(newComment);
            return MapToProtoComment(newComment);
        }

        public override async Task<CommentResponse> UpdateComment(UpdateCommentRequest request, ServerCallContext context)
        {
            var comment = new Comment { Text = request.Text };
            var success = await _commentService.UpdateCommentAsync(request.CommentId, comment);
            
            if (!success) 
                throw new RpcException(new Status(StatusCode.NotFound, "Comment not found"));

            var updated = await _commentService.GetCommentByIdAsync(request.CommentId);
            return MapToProtoComment(updated);
        }

        public override async Task<DeleteCommentResponse> DeleteComment(DeleteCommentRequest request, ServerCallContext context)
        {
            await _commentService.DeleteCommentAsync(request.CommentId);
            return new DeleteCommentResponse { Success = true };
        }

        // --- Mappers ---

        private static PostResponse MapToProtoPost(Post p)
        {
            var resp = new PostResponse
            {
                Id = p.Id ?? "",
                Title = p.Title,
                Description = p.Description,
                LikeCount = p.LikeCount,
                CreatedAt = Timestamp.FromDateTime(p.CreatedAt)
            };
            if (p.ImagePaths != null) resp.ImagePaths.AddRange(p.ImagePaths);
            return resp;
        }

        private static CommentResponse MapToProtoComment(Comment c)
        {
            return new CommentResponse
            {
                Id = c.Id ?? "",
                PostId = c.PostId ?? "",
                AuthorId = c.AuthorId ?? "",
                AuthorName = c.AuthorName ?? "",
                Text = c.Text ?? "",
                CreatedAt = Timestamp.FromDateTime(c.CreatedAt ?? DateTime.UtcNow),
                UpdatedAt = Timestamp.FromDateTime(c.UpdatedAt ?? DateTime.UtcNow)
            };
        }
    }
}