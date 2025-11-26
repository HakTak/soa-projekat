using Blog; // This namespace comes from the Proto generation
using Grpc.Core;
using BLOG.Model;
using BLOG.Services;
using Google.Protobuf.WellKnownTypes;
using System.Linq;
using System.Runtime.CompilerServices;

namespace BLOG.GrpcServices
{
    public class BlogGrpcService: Blog.BlogService.BlogServiceBase
    {
        private readonly Services.BlogService _blogService;
        private readonly CommentService _commentService;

        public BlogGrpcService(Services.BlogService blogService, CommentService commentService)
        {
            _blogService = blogService;
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

        public override async Task<GetAllBlogsResponse> GetAllPosts(Empty request, ServerCallContext context)
        {
            var blogs = await _blogService.GetBlogsAsync();
            var response = new GetAllBlogsResponse();

            foreach (var blog in blogs)
            {
                var blogt = new Blog.Blog
                {
                    Id = blog.Id ?? "",
                    Title = blog.Title,
                    Description = blog.Description,
                    UserName = blog.UserName ?? "",
                    LikeCount = blog.LikeCount,
                    CreatedAt = Timestamp.FromDateTime(blog.CreatedAt.ToUniversalTime())
                };
                response.Blogs.Add(blogt);
            }

            return response;
        }

        public override async Task<BlogResponse> CreatePost(CreateBlogRequest request, ServerCallContext context)
        {
            var newBlog = new Model.Blog
            {
                Title = request.Title,
                Description = request.Description,
                ImagePaths = request.ImagePaths.ToList(),
                CreatedAt = DateTime.UtcNow,
                LikeCount = 0
            };

            await _blogService.CreateBlogAsync(newBlog);
            return MapToProtoBlog(newBlog);
        }

        public override async Task<BlogResponse> ToggleLike(ToggleLikeRequest request, ServerCallContext context)
        {
            var userId = GetUserId(context);
            var updatedBlog = await _blogService.ToggleBlogLikeAsync(request.PostId, userId);

            if (updatedBlog == null)
                throw new RpcException(new Status(StatusCode.NotFound, "Blog not found"));

            return MapToProtoBlog(updatedBlog);
        }

        // ======================
        // COMMENT METHODS
        // ======================
        public override async Task<GetCommentsResponse> GetComments(GetCommentsRequest request, ServerCallContext context)
        {
            var comments = await _commentService.GetCommentsByPostIdAsync(request.PostId);
            var response = new GetCommentsResponse();

            foreach (var comment in comments)
            {
                var newComment = new Blog.Comment
                {
                  PostId = request.PostId,
                  Text = comment.Text,
                  AuthorName = comment.AuthorName,
                  CreatedAt = Timestamp.FromDateTime(comment.CreatedAt.ToUniversalTime()),
                  UpdatedAt = Timestamp.FromDateTime(comment.UpdatedAt.ToUniversalTime())
                };
                response.Comments.Add(newComment);
            }

            return response;
        }

        public override async Task<CommentResponse> CreateComment(CreateCommentRequest request, ServerCallContext context)
        {
            var userId = GetUserId(context);
            var username = context.RequestHeaders.GetValue("x-user-username") ?? "Unknown";

            var newComment = new Model.Comment
            {
                PostId = request.PostId,
                Text = request.Text,
                AuthorName = username,
                CreatedAt = DateTime.UtcNow,
                UpdatedAt = DateTime.UtcNow
            };

            await _commentService.CreateCommentAsync(newComment);
            return MapToProtoComment(newComment);
        }

        public override async Task<CommentResponse> UpdateComment(UpdateCommentRequest request, ServerCallContext context)
        {
            var comment = new Model.Comment
            {
                Id = request.CommentId,
                Text = request.Text
            };

            var success = await _commentService.UpdateCommentAsync(comment);

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

        // ======================
        // Mappers: Domain -> Proto
        // ======================
        private static BlogResponse MapToProtoBlog(Model.Blog blog)
        {
            var resp = new BlogResponse
            {
                Blog = new Blog.Blog
                {
                    Id = blog.Id ?? "",
                    Title = blog.Title,
                    Description = blog.Description,
                    UserName = blog.UserName ?? "",
                    LikeCount = blog.LikeCount,
                    CreatedAt = Timestamp.FromDateTime(blog.CreatedAt.ToUniversalTime())
                }
            };

            if (blog.ImagePaths != null)
                resp.Blog.ImagePaths.AddRange(blog.ImagePaths);

            return resp;
        }

        private static CommentResponse MapToProtoComment(Model.Comment comment)
        {
            return new CommentResponse
            {
                Comment = new Blog.Comment
                {
                    Id = comment.Id ?? "",
                    PostId = comment.PostId ?? "",
                    AuthorName = comment.AuthorName ?? "",
                    Text = comment.Text ?? "",
                    CreatedAt = Timestamp.FromDateTime(comment.CreatedAt.ToUniversalTime()),
                    UpdatedAt = Timestamp.FromDateTime(comment.UpdatedAt.ToUniversalTime())
                }
            };
        }
    }
}